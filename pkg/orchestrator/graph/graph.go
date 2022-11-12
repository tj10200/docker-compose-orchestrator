package graph

import "fmt"

type Edges map[string]Node
type Visited map[string]bool

type Node struct {
	Name    string
	Edges   Edges
	Command string
	File    string
}
type Graph struct {
	Nodes    map[string]Node
	visited  Visited
	topoSort []Node
}

func (g *Graph) AddNode(node Node) {
	if g.Nodes == nil {
		g.Nodes = map[string]Node{}
	}
	if !g.Contains(node.Name) {
		g.Nodes[node.Name] = node
	}
}

func (g *Graph) AddEdge(from, to string) error {
	if g.topoSort != nil {
		g.topoSort = nil
		g.visited = Visited{}
	}

	if !g.Contains(from) {
		return fmt.Errorf("graph doesn't contain from node: %v", from)
	}

	t := g.Nodes[to]
	g.AddNode(t)

	f := g.Nodes[from]
	f.Edges[t.Name] = t
	return nil
}

func (g *Graph) Contains(node string) bool {
	_, ok := g.Nodes[node]
	return ok
}

func (g *Graph) Append(other *Graph) (*Graph, error) {
	for _, n := range other.Nodes {
		g.AddNode(n)
		for _, e := range n.Edges {
			g.AddNode(e)
			if err := g.AddEdge(n.Name, e.Name); err != nil {
				return nil, fmt.Errorf("combine() addEdge failed from(%s) to(%s): err(%v)", n.Name, e.Name, err)
			}
		}
	}
	return g, nil
}

func (g *Graph) TopologySort() []Node {
	if g.topoSort != nil {
		return g.topoSort
	}
	for k, v := range g.Nodes {
		if !g.visited[k] {
			g.topologySort(v)
		}
	}
	return g.topoSort
}

func (g *Graph) topologySort(node Node) {
	if g.visited == nil {
		g.visited = map[string]bool{}
	}
	g.visited[node.Name] = true

	for _, child := range node.Edges {
		if g.visited[child.Name] {
			continue
		}
		g.topologySort(child)
	}
	g.topoSort = append(g.topoSort, node)
}
