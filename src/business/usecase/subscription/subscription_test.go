package subscription

import (
	"context"
	"loverly/lib/appcontext"
	mock_log "loverly/lib/log/mock"
	mock_subscription "loverly/src/business/domain/mock/subscription"
	"loverly/src/business/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock_log.NewMockInterface(ctrl)
	subsMock := mock_subscription.NewMockInterface(ctrl)

	type mockFields struct {
		subsMock *mock_subscription.MockInterface
	}

	mocks := mockFields{
		subsMock: subsMock,
	}

	type args struct {
		ctx context.Context
	}

	allGoods := &entity.Subscription{ID: 1, UserId: 1, Plan: entity.UnlimitedPlan}

	tests := []struct {
		name     string
		mockFunc func(mock mockFields, arg args)
		args     args
		want     *entity.Subscription
		wantErr  bool
	}{
		{
			name: "err get subscription",
			args: args{
				ctx: appcontext.SetUserId(context.Background(), 1),
			},
			want:    &entity.Subscription{},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return(entity.Subscription{}, assert.AnError)
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
				mock.subsMock.EXPECT().GetByUserId(arg.ctx, int64(1)).Return(entity.Subscription{ID: 1, UserId: 1, Plan: entity.UnlimitedPlan}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			d := Init(log, subsMock)
			got, err := d.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock_log.NewMockInterface(ctrl)
	subsMock := mock_subscription.NewMockInterface(ctrl)

	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}

	restoreAll := func() {
		Now = time.Now
	}
	defer restoreAll()

	type mockFields struct {
		subsMock *mock_subscription.MockInterface
	}

	mocks := mockFields{
		subsMock: subsMock,
	}

	type args struct {
		ctx   context.Context
		param entity.SubscriptionParam
	}

	paramMock := entity.SubscriptionParam{Plan: entity.UnlimitedPlan}

	tests := []struct {
		name     string
		mockFunc func(mock mockFields, arg args)
		args     args
		want     error
		wantErr  bool
	}{
		{
			name: "err insert subscription",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    nil,
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().Create(arg.ctx, entity.Subscription{UserId: int64(1), Plan: entity.UnlimitedPlan, StartDate: Now(), EndDate: Now().AddDate(0, 0, 30)}).Return(int64(0), assert.AnError)
			},
		},
		{
			name: "all goods",
			args: args{
				ctx:   appcontext.SetUserId(context.Background(), 1),
				param: paramMock,
			},
			want:    nil,
			wantErr: false,
			mockFunc: func(mock mockFields, arg args) {
				mock.subsMock.EXPECT().Create(arg.ctx, entity.Subscription{UserId: int64(1), Plan: entity.UnlimitedPlan, StartDate: Now(), EndDate: Now().AddDate(0, 0, 30)}).Return(int64(0), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			d := Init(log, subsMock)
			err := d.Create(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// assert.Equal(t, tt.want, err)
		})
	}
}
