package user

import (
	"context"
	"loverly/lib/appcontext"
	mock_log "loverly/lib/log/mock"
	mock_profile "loverly/src/business/domain/mock/profile"
	mock_user "loverly/src/business/domain/mock/user"
	"loverly/src/business/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"

	atomicSQLX "loverly/lib/atomic/sqlx"
)

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock_log.NewMockInterface(ctrl)
	userMock := mock_user.NewMockInterface(ctrl)
	profileMock := mock_profile.NewMockInterface(ctrl)

	tracer := otel.Tracer("test")
	atomicSessionProvider := atomicSQLX.NewSqlxAtomicSessionProvider(nil, tracer, log)

	type mockFields struct {
		userMock *mock_user.MockInterface
	}

	mocks := mockFields{
		userMock: userMock,
	}

	type args struct {
		ctx   context.Context
		param entity.SignInParam
	}

	// allGoods := []entity.ProfileResponse{
	// 	{FullName: "test", Gender: entity.Female, Age: 292},
	// }

	tests := []struct {
		name     string
		mockFunc func(mock mockFields, arg args)
		args     args
		want     *entity.SignInResponse
		wantErr  bool
	}{
		{
			name: "err get user",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: entity.SignInParam{Email: "test", Password: "test"},
			},
			want:    &entity.SignInResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.userMock.EXPECT().GetByEmail(arg.ctx, arg.param.Email).Return(entity.User{}, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			d := Init(log, nil, userMock, profileMock, atomicSessionProvider)
			got, err := d.SignIn(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignIn error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
