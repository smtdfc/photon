package core

import (
	"errors"
	"fmt"
)

type Provider interface {
	provide() any
}

func Provide[T any](target *Module, value T) {
	target.provided = value
}

func Inject(target *Module, edge string, provider Provider) error {
	if target.injected == nil {
		target.injected = make(map[string]any)
	}
	target.injected[edge] = provider.provide()
	return nil
}

func Resolve[T any](target *Module, edge string) (T, error) {
	var zero T
	val, ok := target.injected[edge]
	if !ok {
		return zero, errors.New("Dependency " + edge + " not found")
	}

	casted, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("Cannot assertion for type %T", val)
	}

	return casted, nil
}
