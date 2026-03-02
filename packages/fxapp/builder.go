package fxapp

import (
	"slices"

	"go.uber.org/fx"
)

type builder struct {
	options  []fx.Option
	invokers []fx.Option
}

type Builder interface {
	Option(opts ...fx.Option) Builder
	Invoker(invokers ...interface{}) Builder
	Build() App
}

func NewBuilder() Builder {
	return &builder{
		options: []fx.Option{
			fx.NopLogger,
		},
		invokers: make([]fx.Option, 0),
	}
}

func (b *builder) Option(opts ...fx.Option) Builder {
	b.options = append(b.options, opts...)
	return b
}

func (b *builder) Invoker(invokers ...interface{}) Builder {
	for _, invoker := range invokers {
		b.invokers = append(b.invokers, fx.Invoke(invoker))
	}
	return b
}

func (b *builder) Build() App {
	options := slices.Concat(b.options, b.invokers)
	return NewApp(options...)
}
