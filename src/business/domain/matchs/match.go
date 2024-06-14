package match

import (
	"context"
	"fmt"
	"loverly/lib/atomic"
	"loverly/lib/log"
	"loverly/lib/redis"
	"loverly/src/business/entity"

	atomicSqlx "loverly/lib/atomic/sqlx"
	sqlxUtils "loverly/lib/sqlx"

	"github.com/jmoiron/sqlx"
)

type Interface interface {
	GetByUserId(ctx context.Context, userId int64) ([]entity.Match, error)
	Create(ctx context.Context, param entity.Match) (int64, error)
}

type match struct {
	log               log.Interface
	leaderDB          *sqlx.DB
	followerDB        *sqlx.DB
	rds               redis.Redis
	masterStmts       []*sqlx.Stmt
	slaveStmts        []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

const (
	AllFields = `id, user_id_1, user_id_2, created_at, updated_at, deleted_at`

	GetByUserId = iota

	Create

	GetByUserIddKey = "matchs:getbyuserid:%d"
	DeleteKey       = "matchs:*"
)

var (
	masterQueries = []string{}

	masterNamedQueries = []string{
		Create: `INSERT INTO matchs (user_id_1, user_id_2, created_at, updated_at) 
		VALUES (:user_id_1, :user_id_2, now(), now()) RETURNING id`,
	}

	slaveQueries = []string{
		GetByUserId: fmt.Sprintf("SELECT %s FROM matchs WHERE user_id_1 = $1 OR user_id_2 = $1 AND deleted_at IS NULL", AllFields),
	}
)

func Init(ctx context.Context, log log.Interface, leader *sqlx.DB, follower *sqlx.DB, rds redis.Redis) Interface {
	stmpts, err := sqlxUtils.PrepareQueries(leader, masterQueries)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("PrepareQueries err: %v", err))
		return nil
	}

	namedStmpts, err := sqlxUtils.PrepareNamedQueries(leader, masterNamedQueries)
	if err != nil {
		log.Error(ctx, fmt.Sprintf(")PrepareNamedQueries err: %v", err))
		return nil
	}

	slaveStmpts, err := sqlxUtils.PrepareQueries(follower, slaveQueries)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("PrepareQueries err: %v", err))
		return nil
	}

	return &match{
		log:               log,
		leaderDB:          leader,
		followerDB:        follower,
		rds:               rds,
		masterStmts:       stmpts,
		slaveStmts:        slaveStmpts,
		masterNamedStmpts: namedStmpts,
	}
}

func (m *match) GetByUserId(ctx context.Context, userId int64) ([]entity.Match, error) {
	var matchs []entity.Match

	err := m.rds.WithCache(ctx, fmt.Sprintf(GetByUserIddKey, userId), &matchs, func() (interface{}, error) {
		if err := m.slaveStmts[GetByUserId].SelectContext(ctx, &matchs, userId); err != nil {
			return matchs, err
		}

		return matchs, nil
	})
	if err != nil {
		m.log.Error(ctx, fmt.Sprintf("GetByUserId err: %v", err))
		return matchs, err
	}

	return matchs, nil
}

func (m *match) Create(ctx context.Context, param entity.Match) (int64, error) {
	var matchs entity.Match

	namedStmt, err := m.getNamedStatement(ctx, Create)
	if err != nil {
		m.log.Error(ctx, fmt.Sprintf("getNamedStatement err: %v", err))
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &matchs, param); err != nil {
		m.log.Error(ctx, fmt.Sprintf("CreateMatchs err: %v", err))
		return 0, err
	}

	redisErr := m.rds.DelWithPattern(ctx, DeleteKey)
	if redisErr != nil {
		m.log.Error(ctx, fmt.Sprintf("error when redis delete with pattern: %s, %s", DeleteKey, redisErr))
	}

	return matchs.ID, nil
}

func (m *match) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = m.masterStmts[queryId]
	}
	return statement, err
}

func (m *match) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = m.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}
