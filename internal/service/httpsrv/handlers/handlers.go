package handlers

import (
	"github.com/gophers0/users/internal/repository/postgres"
	gaarx "github.com/zergu1ar/Gaarx"
)

type Handlers struct {
	app *gaarx.App
}

func New(app *gaarx.App) *Handlers {
	return &Handlers{app: app}
}

func (h *Handlers) getDB() *postgres.Repo {
	return h.app.GetDB().(*postgres.Repo)
}
