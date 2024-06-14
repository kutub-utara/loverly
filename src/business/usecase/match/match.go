package match

import (
	"context"
	"fmt"
	"loverly/lib/appcontext"
	"loverly/lib/log"
	match "loverly/src/business/domain/matchs"
	"loverly/src/business/domain/profile"
	"loverly/src/business/entity"
	appErr "loverly/src/errors"
	"slices"
	"time"
)

type Interface interface {
	GetList(ctx context.Context) ([]entity.ProfileResponse, error)
}

type matchs struct {
	log     log.Interface
	match   match.Interface
	profile profile.Interface
}

func Init(log log.Interface, m match.Interface, p profile.Interface) Interface {
	return &matchs{
		log:     log,
		match:   m,
		profile: p,
	}
}

func (m *matchs) GetList(ctx context.Context) ([]entity.ProfileResponse, error) {
	var results []entity.ProfileResponse

	userId := appcontext.GetUserId(ctx)
	if userId < 1 {
		return results, appErr.ErrInvalidUserId
	}

	matchs, err := m.match.GetByUserId(ctx, int64(userId))
	if err != nil {
		return results, err
	}

	var userIds []string
	for _, m := range matchs {
		userId1, userId2 := fmt.Sprintf("%d", m.UserId1), fmt.Sprintf("%d", m.UserId2)
		if !slices.Contains(userIds, userId1) && m.UserId1 != int64(userId) {
			userIds = append(userIds, userId1)
		}

		if !slices.Contains(userIds, userId2) && m.UserId2 != int64(userId) {
			userIds = append(userIds, userId2)
		}
	}

	profiles, err := m.profile.GetByUserIds(ctx, userIds)
	if err != nil {
		return results, err
	}

	for _, p := range profiles {
		days := int(time.Now().Sub(p.BirthDay.Time).Hours() / 24)
		results = append(results, entity.ProfileResponse{
			FullName:  p.FullName,
			Gender:    p.Gender,
			Age:       int64(days / 365),
			Location:  p.Location.String,
			Bio:       p.Bio.String,
			Interest:  p.Interest.String,
			CreatedAt: p.CreatedAt.Time,
		})
	}

	return results, nil
}
