package usecase

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(fx.Annotate(NewRunner, fx.ParamTags(``, `name:"producer"`))),
)
