package parser

import (
	"fmt"
	"github.com/tj10200/docker-compose-orchestrator/pkg/orchestrator/graph"
)

type Service struct {
	Image     string   `yaml:"image"`
	Ports     []string `yaml:"ports"`
	Networks  []string `yaml:"networks"`
	DependsOn []string `yaml:"depends_on"`
	Links     []string `yaml:"links"`
}

type ComposeFile struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

func (c ComposeFile) AsGraph(filePath string) (g graph.Graph, err error) {
	for name := range c.Services {
		g.AddNode(graph.Node{
			Name:  name,
			Edges: graph.Edges{},
			File:  filePath,
		})
	}

	for name, svc := range c.Services {
		for _, dep := range svc.DependsOn {
			if err = g.AddEdge(name, dep); err != nil {
				return g, fmt.Errorf("cannot add link edge from(%s) to(%s): err(%v)", name, dep, err)
			}
		}
		for _, link := range svc.Links {
			if err = g.AddEdge(name, link); err != nil {
				return g, fmt.Errorf("cannot add link edge from(%s) to(%s): err(%v)", name, link, err)
			}
		}
	}

	return
}
