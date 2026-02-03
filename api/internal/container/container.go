package container

import (
	"reflect"
	"sync"
)

type Container struct {
	instances sync.Map
	providers sync.Map
}

type Constructor[T any] func(*Container) (T, error)
type ConfigPath string

func New() *Container {
	return &Container{}
}

func Get[T any](c *Container) (T, error) {
	var zero T

	key := getKey(zero)

	instance, ok := c.instances.Load(key)

	if ok {
		return instance.(T), nil
	}

	provider, ok := c.providers.Load(key)

	if !ok {
		panic("no provider for " + key)
	}

	instance, err := provider.(Constructor[T])(c)
	if err != nil {
		return zero, err
	}
	c.instances.Store(key, instance)
	return instance.(T), nil
}

func Provide[T any](c *Container, constructor Constructor[T]) {
	var zero T
	key := getKey(zero)
	c.providers.Store(key, constructor)
}

func ProvideNamed[T any](c *Container, name string, constructor Constructor[T]) {
	key := getKey(*new(T)) + ":" + name
	c.providers.Store(key, constructor)
}

func GetNamed[T any](c *Container, name string) (T, error) {
	var zero T

	key := getKey(*new(T)) + ":" + name

	if instance, ok := c.instances.Load(key); ok {
		return instance.(T), nil
	}

	provider, exists := c.providers.Load(key)
	if !exists {
		panic("named service not registered: " + key)
	}

	instance, err := provider.(Constructor[T])(c)
	if err != nil {
		return zero, err
	}
	c.instances.Store(key, instance)

	return instance, nil
}

func getKey[T any](_ T) string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}
