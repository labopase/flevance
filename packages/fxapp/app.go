package fxapp

import (
	"context"

	"go.uber.org/fx"
)

type app struct {
	fxApp *fx.App
}

type App interface {
	Run()
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Wait() <-chan fx.ShutdownSignal
}

func NewApp(opts ...fx.Option) App {
	return &app{
		fxApp: fx.New(opts...),
	}
}

func (a *app) Run() {
	a.fxApp.Run()
}

func (a *app) Start(ctx context.Context) error {
	return a.fxApp.Start(ctx)
}

func (a *app) Stop(ctx context.Context) error {
	return a.fxApp.Stop(ctx)
}

func (a *app) Wait() <-chan fx.ShutdownSignal {
	return a.fxApp.Wait()
}
