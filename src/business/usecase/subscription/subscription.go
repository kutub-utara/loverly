package subscription

import (
	"context"
	"database/sql"
	"errors"
	"loverly/lib/appcontext"
	"loverly/lib/log"
	"loverly/src/business/domain/subscription"
	"loverly/src/business/entity"
	"time"

	appErr "loverly/src/errors"
)

var Now = time.Now

type Interface interface {
	Get(ctx context.Context) (*entity.Subscription, error)
	Create(ctx context.Context, param entity.SubscriptionParam) error
}

type subs struct {
	log          log.Interface
	subscription subscription.Interface
}

func Init(log log.Interface, s subscription.Interface) Interface {
	return &subs{
		log:          log,
		subscription: s,
	}
}

func (s *subs) Create(ctx context.Context, param entity.SubscriptionParam) error {
	userId := appcontext.GetUserId(ctx)
	if userId < 1 {
		return appErr.ErrInvalidUserId
	}

	_, err := s.subscription.Create(ctx, entity.Subscription{
		UserId:    int64(userId),
		Plan:      param.Plan,
		StartDate: Now(),
		EndDate:   Now().AddDate(0, 0, 30),
	})

	return err
}

func (s *subs) Get(ctx context.Context) (*entity.Subscription, error) {
	var results entity.Subscription

	userId := appcontext.GetUserId(ctx)
	if userId < 1 {
		return &results, appErr.ErrInvalidUserId
	}

	results, err := s.subscription.GetByUserId(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return &results, err
	}

	return &results, nil
}
