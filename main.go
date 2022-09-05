package main

import (
	"github.com/marianozunino/godollar/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.Module).Run()
}
