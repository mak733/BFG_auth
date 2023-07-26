package controllers

import (
	"BFG_auth/controllers/ubuntu"
	"github.com/pkg/errors"
)

// IdP - Interface for Identity Providers
type Controller interface {
	Name() string
	Time() string
	Disk() string
	Version() string
	Network() string
	Ram() string
}

func NewController(OS string) (Controller, error) {

	switch OS {
	case "Ubuntu":
		return &ubuntu.Ubuntu{}, nil
	default:
		return nil, errors.New("unknown controllers type")
	}
}
