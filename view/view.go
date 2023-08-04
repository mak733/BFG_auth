package view

import (
	"BFG_auth/session"
	"BFG_auth/view/http"
	"github.com/pkg/errors"
)

type View interface {
	StartServer(address string, sessionManager *session.Manager) error
}

func NewView(ViewName string) (View, error) {
	switch ViewName {
	case "http":
		return &http.ViewHttp{}, nil
	default:
		return nil, errors.New("unknown controllers type")
	}
}
