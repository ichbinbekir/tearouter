package router

import (
	"github.com/ichbinbekir/tearouter"
	"github.com/ichbinbekir/tearouter/internal/models/page"
)

func Model() tearouter.Model {
	return tearouter.Model{
		Routes: []tearouter.Route{
			{
				Path:    "/",
				Builder: page.Main,
			},
		},
	}
}
