package controllers

import (
	"BFG_auth/controllers/ubuntu"
	"github.com/pkg/errors"
)

// IdP - Interface for Identity Providers
type API interface {
	Name() string
	Time() string
	Disk() string
	Version() string
	Network() string
	Ram() string
}

func NewController(OS string) (API, error) {

	switch OS {
	case "Ubuntu":
		return &ubuntu.Ubuntu{}, nil
	default:
		return nil, errors.New("unknown controllers type")
	}
}
