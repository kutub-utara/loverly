package user

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
	// Get(ctx context.Context, params entity.user) (entity.user, error)
	// GetById(ctx context.Context, id int64) (entity.user, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Create(ctx context.Context, param entity.User) (int64, error)
}

type user struct {
	log               log.Interface
	leaderDB          *sqlx.DB
	followerDB        *sqlx.DB
	rds               redis.Redis
	masterStmts       []*sqlx.Stmt
	slaveStmts        []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

const (
	AllFields = `id, email, password, verified, created_at, updated_at, deleted_at`

	Get = iota
	// GetById
	GetByEmail

	Create

	// GetListKey    = "users:getlist"
	// GetByIdKey    = "users:getbyid:%d"
	GetByEmailKey = "users:getbyemail:%s"
	DeleteKey     = "users:*"
)

var (
	masterQueries = []string{}

	masterNamedQueries = []string{
		Create: `INSERT INTO users (email, password, created_at, updated_at) 
		VALUES (:email, :password, now(), now()) RETURNING id`,
	}

	slaveQueries = []string{
		Get: fmt.Sprintf("SELECT %s FROM users WHERE deleted_at IS NULL", AllFields),
		// GetById:    fmt.Sprintf("SELECT %s FROM users WHERE id = $1 AND deleted_at IS NULL", AllFields),
		GetByEmail: fmt.Sprintf("SELECT %s FROM users WHERE email = $1 AND deleted_at IS NULL", AllFields),
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

	return &user{
		log:               log,
		leaderDB:          leader,
		followerDB:        follower,
		rds:               rds,
		masterStmts:       stmpts,
		slaveStmts:        slaveStmpts,
		masterNamedStmpts: namedStmpts,
	}
}

func (u *user) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User

	err := u.rds.WithCache(ctx, fmt.Sprintf(GetByEmailKey, email), &user, func() (interface{}, error) {
		if err := u.slaveStmts[GetByEmail].GetContext(ctx, &user, email); err != nil {
			return user, err
		}

		return user, nil
	})
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf("GetByEmail err: %v", err))
		return user, err
	}

	return user, nil
}

func (u *user) Create(ctx context.Context, param entity.User) (int64, error) {
	var user entity.User

	namedStmt, err := u.getNamedStatement(ctx, Create)
	if err != nil {
		u.log.Error(ctx, fmt.Sprintf("getNamedStatement err: %v", err))
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &user, param); err != nil {
		u.log.Error(ctx, fmt.Sprintf("CreateUser err: %v", err))
		return 0, err
	}

	redisErr := u.rds.DelWithPattern(ctx, DeleteKey)
	if redisErr != nil {
		u.log.Error(ctx, fmt.Sprintf("error when redis delete with pattern: %s, %s", DeleteKey, redisErr))
	}

	return user.ID, nil
}

func (r *user) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = r.masterStmts[queryId]
	}
	return statement, err
}

func (r *user) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = r.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}
