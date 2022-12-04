// Copyright 2020 Siigo. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package controller

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"github.com/sirupsen/logrus"
	"sync"
)

// Controller implements the protobuf interface
type Controller struct {
	Mu     *sync.RWMutex
	Bus    cqrs.Dispatcher
	Logger *logrus.Logger
}

// NewController initializes a new Controller struct.
func NewController(bus cqrs.Dispatcher, logger *logrus.Logger) *Controller {
	return &Controller{
		Mu:     &sync.RWMutex{},
		Bus:    bus,
		Logger: logger,
	}
}
