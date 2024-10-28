package model

import (
	"fmt"

	"github.com/rnovatorov/soegur/internal/api/sagaspecpb"
)

type dependencyGraph struct {
	vertices map[string]struct{}
	edges    map[string]map[string]struct{}
}

func buildDependencyGraph(steps []*sagaspecpb.Step) (*dependencyGraph, error) {
	g := &dependencyGraph{
		vertices: make(map[string]struct{}),
		edges:    make(map[string]map[string]struct{}),
	}

	for _, step := range steps {
		if err := g.addVertex(step.GetId()); err != nil {
			return nil, fmt.Errorf("add vertex: %s: %w", step.Id, err)
		}
	}

	for _, step := range steps {
		for _, dep := range step.Dependencies {
			if err := g.addEdge(step.GetId(), dep); err != nil {
				return nil, fmt.Errorf("add edge: %s->%s: %w",
					step.Id, dep, err)
			}
		}
	}

	if g.hasCycles() {
		return nil, ErrDependencyGraphCyclic
	}

	return g, nil
}

func (g *dependencyGraph) addVertex(id string) error {
	if _, ok := g.vertices[id]; ok {
		return ErrDependencyGraphVertexDefinedTwice
	}

	g.vertices[id] = struct{}{}
	g.edges[id] = make(map[string]struct{})

	return nil
}

func (g *dependencyGraph) addEdge(from, to string) error {
	if _, ok := g.vertices[from]; !ok {
		return fmt.Errorf("%w: %s", ErrDependencyGraphVertexNotFound, from)
	}
	if _, ok := g.vertices[to]; !ok {
		return fmt.Errorf("%w: %s", ErrDependencyGraphVertexNotFound, to)
	}

	g.edges[from][to] = struct{}{}

	return nil
}

func (g *dependencyGraph) dependsOn(dependant, dependency string) bool {
	var dfs func(string) bool
	dfs = func(id string) bool {
		for directDependency := range g.edges[id] {
			if directDependency == dependency || dfs(directDependency) {
				return true
			}
		}
		return false
	}

	return dfs(dependant)
}

func (g *dependencyGraph) hasCycles() bool {
	checked := make(map[string]bool, len(g.vertices))
	inStack := make(map[string]bool, len(checked))

	var dfs func(string) bool
	dfs = func(from string) bool {
		if _, ok := checked[from]; ok {
			return false
		}
		defer func() { checked[from] = true }()

		inStack[from] = true
		defer func() { inStack[from] = false }()

		for to := range g.edges[from] {
			if inStack[to] {
				return true
			}
			if dfs(to) {
				return true
			}
		}

		return false
	}

	for id := range g.vertices {
		if dfs(id) {
			return true
		}
	}

	return false
}

func (g *dependencyGraph) transitiveDependencies(id string) []string {
	if _, ok := g.vertices[id]; !ok {
		return nil
	}

	var result []string
	seen := make(map[string]bool)

	var dfs func(string)
	dfs = func(id string) {
		for dep := range g.edges[id] {
			if !seen[dep] {
				seen[dep] = true
				dfs(dep)
				result = append(result, dep)
			}
		}
	}
	dfs(id)

	return result
}
