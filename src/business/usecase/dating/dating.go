package dating

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"loverly/lib/appcontext"
	"loverly/lib/log"
	match "loverly/src/business/domain/matchs"
	"loverly/src/business/domain/profile"
	"loverly/src/business/domain/subscription"
	"loverly/src/business/domain/swipe"
	"loverly/src/business/entity"
	"time"

	appErr "loverly/src/errors"
)

type Interface interface {
	Discovery(ctx context.Context) ([]entity.Discovery, error)
	Swipe(ctx context.Context, param entity.SwipeParam) (entity.SwipeResponse, error)
}

type dating struct {
	log          log.Interface
	subscription subscription.Interface
	profile      profile.Interface
	swipe        swipe.Interface
	match        match.Interface
}

func Init(log log.Interface, subs subscription.Interface, pr profile.Interface, sw swipe.Interface, m match.Interface) Interface {
	return &dating{
		log:          log,
		subscription: subs,
		profile:      pr,
		swipe:        sw,
		match:        m,
	}
}

func (d *dating) Discovery(ctx context.Context) ([]entity.Discovery, error) {
	var results []entity.Discovery

	userId := appcontext.GetUserId(ctx)
	if userId < 1 {
		return results, appErr.ErrInvalidUserId
	}

	access, err := d.checkQuotaLimit(ctx, int64(userId))
	if err != nil {
		return results, err
	}

	// don't have access
	if !access {
		return results, fmt.Errorf("Please subcribe for unlimited discover profile !")
	}

	uProfile, err := d.profile.GetByUserId(ctx, int64(userId))
	if err != nil {
		return results, err
	}

	gender := entity.Male
	if uProfile.Gender == entity.Male {
		gender = entity.Female
	}

	// get profile unique & with opposite gender
	profiles, err := d.profile.GetBySwipe(ctx, int64(userId), gender)
	if err != nil {
		return results, err
	}

	for _, p := range profiles {
		days := int(time.Now().Sub(p.BirthDay.Time).Hours() / 24)
		results = append(results, entity.Discovery{
			ID:       p.ID,
			FullName: p.FullName,
			Age:      int64(days / 365),
			Gender:   p.Gender,
			Bio:      p.Bio.String,
			Location: p.Location.String,
			Interest: p.Interest.String,
		})
	}

	return results, nil
}

func (d *dating) Swipe(ctx context.Context, param entity.SwipeParam) (entity.SwipeResponse, error) {
	var result entity.SwipeResponse

	userId := appcontext.GetUserId(ctx)
	if userId < 1 {
		return result, appErr.ErrInvalidUserId
	}

	access, err := d.checkQuotaLimit(ctx, int64(userId))
	if err != nil {
		return result, err
	}

	// don't have access
	if !access {
		return result, fmt.Errorf("Please subcribe for unlimited discover profile !")
	}

	_, err = d.swipe.Create(ctx, entity.Swipe{
		SwiperId:  int64(userId),
		SwipedId:  param.SwipedId,
		Direction: param.Direction,
	})
	if err != nil {
		return result, err
	}

	// check match or not
	match, err := d.swipe.GetBySwipeId(ctx, param.SwipedId, int64(userId))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return result, err
		}
	}

	if match.ID > 0 && match.Direction == entity.Like && param.Direction == entity.Like {
		_, err = d.match.Create(ctx, entity.Match{
			UserId1: int64(userId),
			UserId2: param.SwipedId,
		})
		if err != nil {
			return result, err
		}

		result.Match = true
	}

	result.Like = false
	if param.Direction == entity.Like {
		result.Like = true
	}

	return result, nil
}

func (d *dating) checkQuotaLimit(ctx context.Context, userId int64) (bool, error) {
	// check quota plan from subscription
	sub, err := d.subscription.GetByPlan(ctx, int64(userId), entity.UnlimitedPlan)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}

	// get swipe today
	swipes, err := d.swipe.GetBySwiperId(ctx, userId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}

	quota := 10                                       // default quota
	if sub.ID == 0 || sub.EndDate.After(time.Now()) { // if don't have a plan or expired plan
		if len(swipes) >= quota {
			return false, fmt.Errorf("Your quotas exceeded!")
		}
	}

	return true, nil
}
