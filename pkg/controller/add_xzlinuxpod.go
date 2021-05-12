package controller

import (
	"github.com/xzlinux/xzlinuxpod-operator/pkg/controller/xzlinuxpod"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, xzlinuxpod.Add)
}
