package profile

import (
	"context"
	"loverly/lib/appcontext"
	mock_log "loverly/lib/log/mock"
	mock_profile "loverly/src/business/domain/mock/profile"
	"loverly/src/business/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock_log.NewMockInterface(ctrl)
	profileMock := mock_profile.NewMockInterface(ctrl)

	type mockFields struct {
		profileMock *mock_profile.MockInterface
	}

	mocks := mockFields{
		profileMock: profileMock,
	}

	type args struct {
		ctx context.Context
	}

	allGoods := entity.ProfileResponse{
		FullName: "test", Gender: entity.Female, Age: 292,
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockFields, arg args)
		args     args
		want     entity.ProfileResponse
		wantErr  bool
	}{
		{
			name: "err get profile",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    entity.ProfileResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.profileMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return(entity.Profile{}, assert.AnError)
			},
		},
		{
			name: "all  goods",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    allGoods,
			wantErr: false,
			mockFunc: func(mock mockFields, arg args) {
				mock.profileMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return(entity.Profile{FullName: "test", Gender: entity.Female}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			d := Init(log, profileMock)
			got, err := d.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ge error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
