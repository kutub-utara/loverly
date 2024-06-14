package atomic

import (
	"context"
	"fmt"
	"loverly/lib/log"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type SqlxAtomicSession struct {
	tx    *sqlx.Tx
	trace trace.Tracer
	log   log.Interface
}

func NewAtomicSession(tx *sqlx.Tx, tr trace.Tracer, log log.Interface) *SqlxAtomicSession {
	return &SqlxAtomicSession{
		tx:    tx,
		trace: tr,
		log:   log,
	}
}

func (s SqlxAtomicSession) Commit(ctx context.Context) error {
	ctx, span := s.trace.Start(ctx, "SqlxAtomicSession/Commit")
	defer span.End()

	err := s.tx.Commit()
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("commit err: %+v", err))
	}
	return err
}

func (s SqlxAtomicSession) Rollback(ctx context.Context) error {
	ctx, span := s.trace.Start(ctx, "SqlxAtomicSession/Rollback")
	defer span.End()

	err := s.tx.Rollback()
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("rollback err: %+v", err))
	}
	return err
}

func (s SqlxAtomicSession) Tx() *sqlx.Tx {
	return s.tx
}
