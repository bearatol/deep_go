package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	types map[string]any
}

func NewContainer() *Container {
	return &Container{
		types: make(map[string]any),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.types[name] = constructor
}

func (c *Container) Resolve(name string) (interface{}, error) {
	typeVal, ok := c.types[name]
	if !ok {
		return nil, errors.New("type is not exist")
	}

	function, ok := typeVal.(func() any)
	if !ok {
		return nil, errors.New("invalid function type")
	}
	return function(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
