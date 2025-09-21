package dix

import (
	"fmt"
)

type DependencyGraph map[string][]*Dependency
type InitOrder []*Dependency

func BuildInitOrder(graph DependencyGraph) (InitOrder, error) {
	result := InitOrder{}
	visit := make(map[string]int)

	const (
		White = 0
		Gray  = 1
		Black = 2
	)

	var dfs func(dep *Dependency) error
	dfs = func(dep *Dependency) error {
		if graph[dep.Name] == nil {
			return fmt.Errorf("cannot resolve dependency %s", dep.Name)
		}
		switch visit[dep.Name] {
		case Gray:
			return fmt.Errorf("circular dependency detected at %s", dep.Name)
		case Black:
			return nil
		}

		visit[dep.Name] = Gray
		for _, child := range graph[dep.Name] {
			if err := dfs(child); err != nil {
				return err
			}
		}
		visit[dep.Name] = Black
		result = append(result, dep)
		return nil
	}

	for name := range graph {
		if visit[name] == White {
			if err := dfs(&Dependency{Name: name}); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func BuildDepGraph(config *Config) DependencyGraph {
	graph := make(DependencyGraph)

	for name, provider := range config.Providers {
		graph[name] = GetDeps(provider)
	}

	return graph
}
