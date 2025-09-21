package dix

import "strings"

type Dependency struct {
	Name      string
	Type      string
	Transient bool
}

func GetDepValue(val string) *Dependency {
	if strings.HasPrefix(val, "inject:") {
		return &Dependency{
			Name:      strings.Split(val, "inject:")[1],
			Type:      "provider",
			Transient: false,
		}
	}

	if strings.HasPrefix(val, "new:") {
		return &Dependency{
			Name:      strings.Split(val, "new:")[1],
			Type:      "provider",
			Transient: true,
		}
	}

	return nil
}

func GetDeps(provider *Provider) []*Dependency {
	deps := []*Dependency{}
	for _, dep := range provider.Dependencies {
		deps = append(deps, GetDepValue(dep))
	}
	return deps
}
