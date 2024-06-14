package atomic

import (
	"context"
	"fmt"
	"loverly/lib/atomic"
	"loverly/lib/log"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type SqlxAtomicSessionProvider struct {
	db    *sqlx.DB
	trace trace.Tracer
	log   log.Interface
}

func NewSqlxAtomicSessionProvider(db *sqlx.DB, tr trace.Tracer, log log.Interface) *SqlxAtomicSessionProvider {
	return &SqlxAtomicSessionProvider{
		db:    db,
		trace: tr,
		log:   log,
	}
}

func (r *SqlxAtomicSessionProvider) BeginSession(ctx context.Context) (*atomic.AtomicSessionContext, error) {
	ctx, span := r.trace.Start(ctx, "SqlxAtomicSessionProvider/BeginSession")
	defer span.End()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf("begin tx err: %+v", err))
		return nil, err
	}

	atomicSession := NewAtomicSession(tx, r.trace, r.log)
	return atomic.NewAtomicSessionContext(ctx, atomicSession), nil
}
