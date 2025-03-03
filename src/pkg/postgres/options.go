package postgres

import "time"

type Option func(*Postgres)

func MaxPullSize(size int) Option {
	return func(p *Postgres) {
		p.maxPoolsize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(p *Postgres) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgres) {
		p.connTimeout = timeout
	}
}
