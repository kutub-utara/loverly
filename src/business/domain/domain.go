package domain

import (
	"context"
	"loverly/lib/log"
	"loverly/lib/redis"
	match "loverly/src/business/domain/matchs"
	"loverly/src/business/domain/profile"
	"loverly/src/business/domain/subscription"
	"loverly/src/business/domain/swipe"
	"loverly/src/business/domain/user"
	"loverly/src/config"

	"github.com/jmoiron/sqlx"
)

type Domains struct {
	User         user.Interface
	Subscription subscription.Interface
	Swipe        swipe.Interface
	Profile      profile.Interface
	Match        match.Interface
}

type InitParam struct {
	Log        log.Interface
	LeaderDB   *sqlx.DB
	FollowerDB *sqlx.DB
	Rds        redis.Redis
	Cfg        *config.Configuration
}

func Init(ctx context.Context, params InitParam) *Domains {
	return &Domains{
		User:         user.Init(ctx, params.Log, params.LeaderDB, params.FollowerDB, params.Rds),
		Subscription: subscription.Init(ctx, params.Log, params.LeaderDB, params.FollowerDB, params.Rds),
		Swipe:        swipe.Init(ctx, params.Log, params.LeaderDB, params.FollowerDB, params.Rds),
		Profile:      profile.Init(ctx, params.Log, params.LeaderDB, params.FollowerDB, params.Rds),
		Match:        match.Init(ctx, params.Log, params.LeaderDB, params.FollowerDB, params.Rds),
	}
}
