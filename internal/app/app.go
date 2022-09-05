package app

import (
	"context"
	"os"

	"github.com/marianozunino/godollar/internal/app/database"
	"github.com/marianozunino/godollar/internal/app/repository"
	"github.com/marianozunino/godollar/internal/app/route"
	"github.com/marianozunino/godollar/internal/app/server"
	"github.com/marianozunino/godollar/internal/app/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(registerHooks),
	server.Module,
	route.Module,
	service.Module,
	database.Module,
	repository.Module,
)

func registerHooks(lifecycle fx.Lifecycle,
	instance *server.GinHandler) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				// read port from env "PORT"
				// if not set, use default port 8080
				var port = os.Getenv("PORT")
				if port == "" {
					port = "5000"
				}
				go instance.Gin.Run(":" + port)
				return nil
			},
		},
	)
}
