package subscription

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
	GetByUserId(ctx context.Context, userId int64) (entity.Subscription, error)
	GetByPlan(ctx context.Context, userId int64, plan string) (entity.Subscription, error)
	Create(ctx context.Context, param entity.Subscription) (int64, error)
}

type subs struct {
	log               log.Interface
	leaderDB          *sqlx.DB
	followerDB        *sqlx.DB
	rds               redis.Redis
	masterStmts       []*sqlx.Stmt
	slaveStmts        []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

const (
	AllFields = `id, user_id, plan, start_date, end_date, created_at, updated_at, deleted_at`

	GetByUserId = iota
	GetByPlan

	Create

	GetByUserIdKey = "subscriptions:getbyuserid:%d"
	GetByPlanKey   = "subscriptions:getbyplan:%d:%s"
	DeleteKey      = "subscriptions:*"
)

var (
	masterQueries = []string{}

	masterNamedQueries = []string{
		Create: `INSERT INTO subscriptions (user_id, plan, start_date, end_date, created_at, updated_at) 
		VALUES (:user_id, :plan, :start_date, :end_date, now(), now()) RETURNING id`,
	}

	slaveQueries = []string{
		GetByUserId: fmt.Sprintf("SELECT %s FROM subscriptions WHERE user_id = $1 AND deleted_at IS NULL", AllFields),
		GetByPlan:   fmt.Sprintf("SELECT %s FROM subscriptions WHERE user_id = $1 AND plan = $2 AND deleted_at IS NULL", AllFields),
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

	return &subs{
		log:               log,
		leaderDB:          leader,
		followerDB:        follower,
		rds:               rds,
		masterStmts:       stmpts,
		slaveStmts:        slaveStmpts,
		masterNamedStmpts: namedStmpts,
	}
}

func (s *subs) GetByUserId(ctx context.Context, userId int64) (entity.Subscription, error) {
	var subscription entity.Subscription

	err := s.rds.WithCache(ctx, fmt.Sprintf(GetByUserIdKey, userId), &subscription, func() (interface{}, error) {
		if err := s.slaveStmts[GetByUserId].GetContext(ctx, &subscription, userId); err != nil {
			return subscription, err
		}

		return subscription, nil
	})
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("GetByUserId err: %v", err))
		return subscription, err
	}

	return subscription, nil
}

func (s *subs) GetByPlan(ctx context.Context, userId int64, plan string) (entity.Subscription, error) {
	var subscription entity.Subscription

	err := s.rds.WithCache(ctx, fmt.Sprintf(GetByPlanKey, userId, plan), &subscription, func() (interface{}, error) {
		if err := s.slaveStmts[GetByPlan].GetContext(ctx, &subscription, userId, plan); err != nil {
			return subscription, err
		}

		return subscription, nil
	})
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("GetByPlan err: %v", err))
		return subscription, err
	}

	return subscription, nil
}

func (s *subs) Create(ctx context.Context, param entity.Subscription) (int64, error) {
	var user entity.User

	namedStmt, err := s.getNamedStatement(ctx, Create)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("getNamedStatement err: %v", err))
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &user, param); err != nil {
		s.log.Error(ctx, fmt.Sprintf("CreateSubscription err: %v", err))
		return 0, err
	}

	redisErr := s.rds.DelWithPattern(ctx, DeleteKey)
	if redisErr != nil {
		s.log.Error(ctx, fmt.Sprintf("error when redis delete with pattern: %s, %s", DeleteKey, redisErr))
	}

	return user.ID, nil
}

func (s *subs) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = s.masterStmts[queryId]
	}
	return statement, err
}

func (s *subs) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = s.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}
