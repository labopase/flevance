package fxapp

import "go.uber.org/fx"

type builder struct {
	options  []fx.Option
	invokers []fx.Option
}
