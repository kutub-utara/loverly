package atomic

import (
	"context"
	"errors"
	"fmt"

	"loverly/lib/log"
)

var InvalidAtomicSessionProvider error = errors.New("invalid_atomic_session_provider")

type AtomicSessionProvider interface {
	BeginSession(ctx context.Context) (*AtomicSessionContext, error)
}

type AtomicSession interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type AtomicSessionContext struct {
	context.Context
	AtomicSession
}

func NewAtomicSessionContext(ctx context.Context, session AtomicSession) *AtomicSessionContext {
	return &AtomicSessionContext{
		Context:       ctx,
		AtomicSession: session,
	}
}

func Atomic(ctx context.Context, provider AtomicSessionProvider, log log.Interface, fn func(ctx context.Context) error) error {
	sessionCtx, err := provider.BeginSession(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if v := recover(); v != nil {
			if rbErr := sessionCtx.Rollback(ctx); rbErr != nil {
				log.Error(ctx, fmt.Sprintf("rollback from recover err %+v", err))
			}
			panic(v)
		}
	}()

	if err := fn(sessionCtx); err != nil {
		log.Error(ctx, fmt.Sprintf("atomic function return err, rollingback. err: %+v", err))
		if rbErr := sessionCtx.Rollback(ctx); rbErr != nil {
			log.Error(ctx, fmt.Sprintf("rollback err: %+v", err))
		}
		return err
	}

	if cmErr := sessionCtx.Commit(ctx); cmErr != nil {
		log.Error(ctx, fmt.Sprintf("commit err: %+v", err))
		return cmErr
	}

	return nil
}
