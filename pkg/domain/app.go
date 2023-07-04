package domain

import "context"

type App interface {
	Init() error
	Run(ctx context.Context) error
}
