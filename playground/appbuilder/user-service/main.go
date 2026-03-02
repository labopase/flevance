package main

import (
	"github.com/labopase/flevance/packages/fxapp"
	module "github.com/labopase/flevance/playground/appbuilder"
)

func main() {
	builder := fxapp.NewBuilder()

	builder.Option(
		module.AppModule,
		module.LoggerModule,
		module.EchoModule,
	)

	app := builder.Build()
	app.Run()
}
