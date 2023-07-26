package view

import (
	"BFG_auth/controllers"
	"github.com/pkg/errors"
)

type View interface {
	StartServer() bool
}

func NewView(ViewName string, controller controllers.Controller) (View, error) {

	switch ViewName {

	default:
		return nil, errors.New("unknown controllers type")
	}
}
