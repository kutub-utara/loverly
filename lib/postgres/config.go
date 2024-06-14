package postgres

import "time"

type PostgresConfig struct {
	ConnectionUrl      string
	MaxPoolSize        int
	MaxIdleConnections int
	ConnMaxIdleTime    time.Duration
	ConnMaxLifeTime    time.Duration
}
