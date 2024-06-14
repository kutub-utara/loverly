package user

import (
	"context"
	"loverly/lib/atomic"
	"loverly/lib/jwt"
	"loverly/lib/log"
	"loverly/src/business/domain/profile"
	"loverly/src/business/domain/user"
	"loverly/src/business/entity"

	"golang.org/x/crypto/bcrypt"
)

type Interface interface {
	SignIn(ctx context.Context, params entity.SignInParam) (*entity.SignInResponse, error)
	SignUp(ctx context.Context, params entity.SignUpParam) (*entity.SignUpResponse, error)
}

type customer struct {
	log     log.Interface
	user    user.Interface
	profile profile.Interface
	jwt     *jwt.TokenProvider
	atomic  atomic.AtomicSessionProvider
}

func Init(log log.Interface, jwt *jwt.TokenProvider, u user.Interface, p profile.Interface, a atomic.AtomicSessionProvider) Interface {
	return &customer{
		log:     log,
		user:    u,
		profile: p,
		jwt:     jwt,
		atomic:  a,
	}
}

func (c *customer) SignIn(ctx context.Context, params entity.SignInParam) (*entity.SignInResponse, error) {
	resp := &entity.SignInResponse{}

	user, err := c.user.GetByEmail(ctx, params.Email)
	if err != nil {
		return resp, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		return resp, err
	}

	token, err := c.jwt.NewAccessToken(ctx, user.ID, []string{}, jwt.AccessTypeOnline)
	if err != nil {
		return resp, err
	}

	resp = &entity.SignInResponse{
		ID:      user.ID,
		Email:   user.Email,
		Verifed: user.Verifed,
		Token:   token,
	}

	return resp, nil
}

func (c *customer) SignUp(ctx context.Context, params entity.SignUpParam) (*entity.SignUpResponse, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return &entity.SignUpResponse{}, err
	}

	err = atomic.Atomic(ctx, c.atomic, c.log, func(ctx context.Context) error {
		userId, err := c.user.Create(ctx, entity.User{
			Email:    params.Email,
			Password: string(hashPassword),
		})
		if err != nil {
			return err
		}

		_, err = c.profile.Create(ctx, entity.Profile{
			UserId:   userId,
			FullName: params.FullName,
			Gender:   params.Gender,
		})

		return err
	})
	if err != nil {
		return &entity.SignUpResponse{}, err
	}

	return &entity.SignUpResponse{
		NextState: entity.NextStateLogin,
	}, nil
}
