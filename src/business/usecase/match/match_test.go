package match

import (
	"context"
	"loverly/lib/appcontext"
	mock_log "loverly/lib/log/mock"
	mock_match "loverly/src/business/domain/mock/match"
	mock_profile "loverly/src/business/domain/mock/profile"
	"loverly/src/business/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock_log.NewMockInterface(ctrl)
	profileMock := mock_profile.NewMockInterface(ctrl)
	matchMock := mock_match.NewMockInterface(ctrl)

	type mockFields struct {
		profileMock *mock_profile.MockInterface
		matchMock   *mock_match.MockInterface
	}

	mocks := mockFields{
		profileMock: profileMock,
		matchMock:   matchMock,
	}

	type args struct {
		ctx context.Context
	}

	allGoods := []entity.ProfileResponse{
		{FullName: "test", Gender: entity.Female, Age: 292},
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockFields, arg args)
		args     args
		want     []entity.ProfileResponse
		wantErr  bool
	}{
		{
			name: "err get match",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    nil,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.matchMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return([]entity.Match{}, assert.AnError)
			},
		},
		{
			name: "err get profile by user ids",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    nil,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.matchMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return([]entity.Match{{UserId1: 1, UserId2: 2}}, nil)
				mock.profileMock.EXPECT().GetByUserIds(arg.ctx, []string{"2"}).Return([]entity.Profile{}, assert.AnError)
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
				mock.matchMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return([]entity.Match{{UserId1: 1, UserId2: 2}}, nil)
				mock.profileMock.EXPECT().GetByUserIds(arg.ctx, []string{"2"}).Return([]entity.Profile{{FullName: "test", Gender: entity.Female}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			d := Init(log, matchMock, profileMock)
			got, err := d.GetList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetList error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
