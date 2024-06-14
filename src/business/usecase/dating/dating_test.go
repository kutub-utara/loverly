package dating

import (
	"context"
	"loverly/lib/appcontext"
	mock_log "loverly/lib/log/mock"
	mock_match "loverly/src/business/domain/mock/match"
	mock_profile "loverly/src/business/domain/mock/profile"
	mock_subscription "loverly/src/business/domain/mock/subscription"
	mock_swipe "loverly/src/business/domain/mock/swipe"
	"loverly/src/business/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDiscovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock_log.NewMockInterface(ctrl)
	subsMock := mock_subscription.NewMockInterface(ctrl)
	profileMock := mock_profile.NewMockInterface(ctrl)
	swipeMock := mock_swipe.NewMockInterface(ctrl)
	matchMock := mock_match.NewMockInterface(ctrl)

	type mockFields struct {
		subsMock    *mock_subscription.MockInterface
		profileMock *mock_profile.MockInterface
		swipeMock   *mock_swipe.MockInterface
		matchMock   *mock_match.MockInterface
	}

	mocks := mockFields{
		subsMock:    subsMock,
		profileMock: profileMock,
		swipeMock:   swipeMock,
		matchMock:   matchMock,
	}

	type args struct {
		ctx context.Context
	}

	swipesMax := []entity.Swipe{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10}}
	swipesMin := []entity.Swipe{{ID: 1}}
	allGoods := []entity.Discovery{
		{FullName: "test", Gender: entity.Female, Age: 292},
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockFields, arg args)
		args     args
		want     []entity.Discovery
		wantErr  bool
	}{
		{
			name: "err get subscription",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    nil,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, assert.AnError)
			},
		},
		{
			name: "err get swipe",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    nil,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return([]entity.Swipe{}, assert.AnError)
			},
		},
		{
			name: "err quota exceeded",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    nil,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMax, nil)
			},
		},
		{
			name: "err get profile",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    nil,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMin, nil)
				mock.profileMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return(entity.Profile{}, assert.AnError)
			},
		},
		{
			name: "err get profile for match",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    nil,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMin, nil)
				mock.profileMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return(entity.Profile{Gender: entity.Male}, nil)
				mock.profileMock.EXPECT().GetBySwipe(arg.ctx, int64(1), entity.Female).Return([]entity.Profile{}, assert.AnError)
			},
		},
		{
			name: "all goods",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    allGoods,
			wantErr: false,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMin, nil)
				mock.profileMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return(entity.Profile{Gender: entity.Male}, nil)
				mock.profileMock.EXPECT().GetBySwipe(arg.ctx, int64(1), entity.Female).Return([]entity.Profile{{FullName: "test", Gender: entity.Female}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			d := Init(log, subsMock, profileMock, swipeMock, matchMock)
			got, err := d.Discovery(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Discover error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSwipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock_log.NewMockInterface(ctrl)
	subsMock := mock_subscription.NewMockInterface(ctrl)
	profileMock := mock_profile.NewMockInterface(ctrl)
	swipeMock := mock_swipe.NewMockInterface(ctrl)
	matchMock := mock_match.NewMockInterface(ctrl)

	type mockFields struct {
		subsMock    *mock_subscription.MockInterface
		profileMock *mock_profile.MockInterface
		swipeMock   *mock_swipe.MockInterface
		matchMock   *mock_match.MockInterface
	}

	mocks := mockFields{
		subsMock:    subsMock,
		profileMock: profileMock,
		swipeMock:   swipeMock,
		matchMock:   matchMock,
	}

	type args struct {
		ctx   context.Context
		param entity.SwipeParam
	}

	swipesMax := []entity.Swipe{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10}}
	swipesMin := []entity.Swipe{{ID: 1}}
	allGoods := entity.SwipeResponse{
		Like:  true,
		Match: true,
	}
	resp := entity.SwipeResponse{}
	paramMock := entity.SwipeParam{SwipedId: 2, Direction: entity.Like}

	tests := []struct {
		name     string
		mockFunc func(mock mockFields, arg args)
		args     args
		want     entity.SwipeResponse
		wantErr  bool
	}{
		{
			name: "err get subscription",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    resp,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, assert.AnError)
			},
		},
		{
			name: "err get swipe",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    resp,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return([]entity.Swipe{}, assert.AnError)
			},
		},
		{
			name: "err quota exceeded",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    resp,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMax, nil)
			},
		},
		{
			name: "err insert swipe",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    resp,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMin, nil)
				mock.swipeMock.EXPECT().Create(arg.ctx, entity.Swipe{SwiperId: int64(1), SwipedId: arg.param.SwipedId, Direction: arg.param.Direction}).Return(int64(9), assert.AnError)
			},
		},
		{
			name: "err get by swipe id",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    resp,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMin, nil)
				mock.swipeMock.EXPECT().Create(arg.ctx, entity.Swipe{SwiperId: int64(1), SwipedId: arg.param.SwipedId, Direction: arg.param.Direction}).Return(int64(1), nil)
				mock.swipeMock.EXPECT().GetBySwipeId(arg.ctx, arg.param.SwipedId, int64(1)).Return(entity.Swipe{}, assert.AnError)
			},
		},
		{
			name: "err insert match",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    resp,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMin, nil)
				mock.swipeMock.EXPECT().Create(arg.ctx, entity.Swipe{SwiperId: int64(1), SwipedId: arg.param.SwipedId, Direction: arg.param.Direction}).Return(int64(1), nil)
				mock.swipeMock.EXPECT().GetBySwipeId(arg.ctx, arg.param.SwipedId, int64(1)).Return(entity.Swipe{ID: 2, Direction: entity.Like}, nil)
				mock.matchMock.EXPECT().Create(arg.ctx, entity.Match{UserId1: int64(1), UserId2: int64(2)}).Return(int64(0), assert.AnError)
			},
		},
		{
			name: "all goods",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    allGoods,
			wantErr: false,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByPlan(arg.ctx, int64(1), entity.UnlimitedPlan).Return(entity.Subscription{}, nil)
				mock.swipeMock.EXPECT().GetBySwiperId(arg.ctx, int64(1)).Return(swipesMin, nil)
				mock.swipeMock.EXPECT().Create(arg.ctx, entity.Swipe{SwiperId: int64(1), SwipedId: arg.param.SwipedId, Direction: arg.param.Direction}).Return(int64(1), nil)
				mock.swipeMock.EXPECT().GetBySwipeId(arg.ctx, arg.param.SwipedId, int64(1)).Return(entity.Swipe{ID: 2, Direction: entity.Like}, nil)
				mock.matchMock.EXPECT().Create(arg.ctx, entity.Match{UserId1: int64(1), UserId2: int64(2)}).Return(int64(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			d := Init(log, subsMock, profileMock, swipeMock, matchMock)
			got, err := d.Swipe(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Swipe error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
