package core

import (
	"log"
)

type Provider interface {
	GetProvide() any
}

func Inject[T any](target *Module, edge string) T {
	value := target.GetInject(edge)
	if value == nil {
		log.Fatalf("Dependency not found: %q in module %s", edge, target.Name)
	}
	v, ok := value.(T)
	if !ok {
		log.Fatalf("Type mismatch: expected %T but got %T", *new(T), value)
	}
	return v
}

func InjectTo(edge string, target *Module, provider Provider) {
	target.Inject(edge, provider)
}
