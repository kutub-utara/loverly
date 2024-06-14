package profile

import (
	"context"
	"fmt"
	"loverly/lib/atomic"
	"loverly/lib/log"
	"loverly/lib/redis"
	"loverly/src/business/entity"
	"strconv"
	"strings"

	atomicSqlx "loverly/lib/atomic/sqlx"
	sqlxUtils "loverly/lib/sqlx"

	"github.com/jmoiron/sqlx"
)

type Interface interface {
	GetByUserId(ctx context.Context, userId int64) (entity.Profile, error)
	GetByUserIds(ctx context.Context, userId []string) ([]entity.Profile, error)
	GetBySwipe(ctx context.Context, userId int64, gender string) ([]entity.Profile, error)
	Create(ctx context.Context, param entity.Profile) (int64, error)
}

type profile struct {
	log               log.Interface
	leaderDB          *sqlx.DB
	followerDB        *sqlx.DB
	rds               redis.Redis
	masterStmts       []*sqlx.Stmt
	slaveStmts        []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

const (
	AllFields = `id, user_id, name, birthday, gender, location, bio, profile_picture, interests, created_at, updated_at, deleted_at`

	GetBySwipe = iota
	GetByUserId
	GetByUserIds

	Create

	GetBySwipedKey  = "profiles:getbyswipe:%d:%s"
	GetByUserIdKey  = "profiles:getbyuserid:%d"
	GetByUserIdsKey = "profiles:getbyuserids:%s"
	DeleteKey       = "profiles:*"
)

var (
	masterQueries = []string{}

	masterNamedQueries = []string{
		Create: `INSERT INTO profiles (user_id, name, birthday, gender, location, bio, profile_picture, interests, created_at, updated_at) 
		VALUES (:user_id, :name, :birthday, :gender, :location, :bio, :profile_picture, :interests, now(), now()) RETURNING id`,
	}

	slaveQueries = []string{
		GetBySwipe:   fmt.Sprintf("SELECT %s FROM profiles WHERE user_id NOT IN (SELECT swiped_id FROM swipes WHERE swiper_id = $1 and DATE(created_at) = CURRENT_DATE) AND gender = $2 AND deleted_at IS NULL", AllFields),
		GetByUserId:  fmt.Sprintf("SELECT %s FROM profiles WHERE user_id = $1 AND deleted_at IS NULL", AllFields),
		GetByUserIds: fmt.Sprintf("SELECT %s FROM profiles WHERE user_id = ANY($1) AND deleted_at IS NULL", AllFields),
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

	return &profile{
		log:               log,
		leaderDB:          leader,
		followerDB:        follower,
		rds:               rds,
		masterStmts:       stmpts,
		slaveStmts:        slaveStmpts,
		masterNamedStmpts: namedStmpts,
	}
}

func (p *profile) GetByUserId(ctx context.Context, userId int64) (entity.Profile, error) {
	var profile entity.Profile

	err := p.rds.WithCache(ctx, fmt.Sprintf(GetByUserIdKey, userId), &profile, func() (interface{}, error) {
		if err := p.slaveStmts[GetByUserId].GetContext(ctx, &profile, userId); err != nil {
			return profile, err
		}

		return profile, nil
	})
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("GetByuserId err: %v", err))
		return profile, err
	}

	return profile, nil
}

func (p *profile) GetByUserIds(ctx context.Context, userId []string) ([]entity.Profile, error) {
	var profiles []entity.Profile

	userIds := fmt.Sprintf("{%s}", strings.Join(userId, ","))
	err := p.rds.WithCache(ctx, fmt.Sprintf(GetByUserIdsKey, userIds), &profiles, func() (interface{}, error) {
		if err := p.slaveStmts[GetByUserIds].SelectContext(ctx, &profiles, userIds); err != nil {
			return profiles, err
		}

		return profiles, nil
	})
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("GetByUserIds err: %v", err))
		return profiles, err
	}

	return profiles, nil
}

func (p *profile) GetBySwipe(ctx context.Context, userId int64, gender string) ([]entity.Profile, error) {
	var profiles []entity.Profile

	err := p.rds.WithCache(ctx, fmt.Sprintf(GetBySwipedKey, userId, gender), &profiles, func() (interface{}, error) {
		if err := p.slaveStmts[GetBySwipe].SelectContext(ctx, &profiles, userId, gender); err != nil {
			return profiles, err
		}

		return profiles, nil
	})
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("GetBySwipe err: %v", err))
		return profiles, err
	}

	return profiles, nil
}

func (p *profile) Create(ctx context.Context, param entity.Profile) (int64, error) {
	var profile entity.Profile

	namedStmt, err := p.getNamedStatement(ctx, Create)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("getNamedStatement err: %v", err))
		return 0, err
	}

	if err = namedStmt.GetContext(ctx, &profile, param); err != nil {
		p.log.Error(ctx, fmt.Sprintf("CreateProfile err: %v", err))
		return 0, err
	}

	redisErr := p.rds.DelWithPattern(ctx, DeleteKey)
	if redisErr != nil {
		p.log.Error(ctx, fmt.Sprintf("error when redis delete with pattern: %s, %s", DeleteKey, redisErr))
	}

	return profile.ID, nil
}

func (p *profile) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = p.masterStmts[queryId]
	}
	return statement, err
}

func (p *profile) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = p.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}

func int64SliceToString(slice []int64) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(strSlice, ",")
}
