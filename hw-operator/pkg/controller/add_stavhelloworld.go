package controller

import (
	"github.com/stavco9/go-test/pkg/controller/stavhelloworld"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, stavhelloworld.Add)
}
