package swipe

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
	GetBySwiperId(ctx context.Context, swiperId int64) ([]entity.Swipe, error)
	GetBySwipeId(ctx context.Context, swiperId, swipedId int64) (entity.Swipe, error)
	Create(ctx context.Context, param entity.Swipe) (int64, error)
}

type swipe struct {
	log               log.Interface
	leaderDB          *sqlx.DB
	followerDB        *sqlx.DB
	rds               redis.Redis
	masterStmts       []*sqlx.Stmt
	slaveStmts        []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

const (
	AllFields = `id, swiper_id, swiped_id, direction, created_at, updated_at, deleted_at`

	GetBySwiperId = iota
	GetBySwipeId

	Create

	GetBySwipeIdKey  = "swipes:getbyswipeid:%d:%d"
	GetBySwiperIdKey = "swipes:getbyswiperid:%d"
	DeleteKey        = "swipes:*"
)

var (
	masterQueries = []string{}

	masterNamedQueries = []string{
		Create: `INSERT INTO swipes (swiper_id, swiped_id, direction, created_at, updated_at) 
		VALUES (:swiper_id, :swiped_id, :direction, now(), now()) RETURNING id`,
	}

	slaveQueries = []string{
		GetBySwiperId: fmt.Sprintf("SELECT %s FROM swipes WHERE swiper_id = $1 AND DATE(created_at) = CURRENT_DATE AND deleted_at IS NULL", AllFields),
		GetBySwipeId:  fmt.Sprintf("SELECT %s FROM swipes WHERE swiper_id = $1 AND swiped_id = $2 AND deleted_at IS NULL", AllFields),
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

	return &swipe{
		log:               log,
		leaderDB:          leader,
		followerDB:        follower,
		rds:               rds,
		masterStmts:       stmpts,
		slaveStmts:        slaveStmpts,
		masterNamedStmpts: namedStmpts,
	}
}

func (s *swipe) GetBySwiperId(ctx context.Context, swiperId int64) ([]entity.Swipe, error) {
	var swipes []entity.Swipe

	err := s.rds.WithCache(ctx, fmt.Sprintf(GetBySwiperIdKey, swiperId), &swipes, func() (interface{}, error) {
		if err := s.slaveStmts[GetBySwiperId].SelectContext(ctx, &swipes, swiperId); err != nil {
			return swipes, err
		}

		return swipes, nil
	})
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("GetBySwiperId err: %v", err))
		return swipes, err
	}

	return swipes, nil
}

func (s *swipe) GetBySwipeId(ctx context.Context, swiperId, swipedId int64) (entity.Swipe, error) {
	var swipes entity.Swipe

	err := s.rds.WithCache(ctx, fmt.Sprintf(GetBySwipeIdKey, swiperId, swipedId), &swipes, func() (interface{}, error) {
		if err := s.slaveStmts[GetBySwipeId].GetContext(ctx, &swipes, swiperId, swipedId); err != nil {
			return swipes, err
		}

		return swipes, nil
	})
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("GetBySwipeId err: %v", err))
		return swipes, err
	}

	return swipes, nil
}

func (s *swipe) Create(ctx context.Context, param entity.Swipe) (int64, error) {
	var swipes entity.Swipe

	namedStmt, err := s.getNamedStatement(ctx, Create)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("getNamedStatement err: %v", err))
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &swipes, param); err != nil {
		s.log.Error(ctx, fmt.Sprintf("CreateSwipes err: %v", err))
		return 0, err
	}

	redisErr := s.rds.DelWithPattern(ctx, DeleteKey)
	if redisErr != nil {
		s.log.Error(ctx, fmt.Sprintf("error when redis delete with pattern: %s, %s", DeleteKey, redisErr))
	}

	redisErr = s.rds.DelWithPattern(ctx, "profiles:*")
	if redisErr != nil {
		s.log.Error(ctx, fmt.Sprintf("error when redis delete with pattern: %s, %s", DeleteKey, redisErr))
	}

	return swipes.ID, nil
}

func (s *swipe) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (s *swipe) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
