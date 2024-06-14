package profile

import (
	"context"
	"loverly/lib/appcontext"
	"loverly/lib/log"
	"loverly/src/business/domain/profile"
	"loverly/src/business/entity"
	appErr "loverly/src/errors"
	"time"
)

type Interface interface {
	Get(ctx context.Context) (entity.ProfileResponse, error)
}

type profiles struct {
	log     log.Interface
	profile profile.Interface
}

func Init(log log.Interface, p profile.Interface) Interface {
	return &profiles{
		log:     log,
		profile: p,
	}
}

func (p *profiles) Get(ctx context.Context) (entity.ProfileResponse, error) {
	var results entity.ProfileResponse

	userId := appcontext.GetUserId(ctx)
	if userId < 1 {
		return results, appErr.ErrInvalidUserId
	}

	pf, err := p.profile.GetByUserId(ctx, int64(userId))
	if err != nil {
		return results, err
	}

	days := int(time.Now().Sub(pf.BirthDay.Time).Hours() / 24)
	results = entity.ProfileResponse{
		FullName:  pf.FullName,
		Gender:    pf.Gender,
		Age:       int64(days / 365),
		Location:  pf.Location.String,
		Bio:       pf.Bio.String,
		Interest:  pf.Interest.String,
		CreatedAt: pf.CreatedAt.Time,
	}

	return results, nil
}
