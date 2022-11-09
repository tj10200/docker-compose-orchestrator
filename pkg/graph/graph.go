package graph

type Edges map[string]Node
type Visited map[string]bool

type Node struct {
	Name    string
	Edges   Edges
	Command string
}
type Graph struct {
	Nodes    map[string]Node
	visited  Visited
	topoSort []Node
}

func (g *Graph) AddNode(node Node) {
	if !g.Contains(node.Name) {
		g.Nodes[node.Name] = node
	}
}

func (g *Graph) AddEdge(from, to Node) {
	if g.topoSort != nil {
		g.topoSort = nil
		g.visited = Visited{}
	}
	g.AddNode(from)
	g.AddNode(to)

	from.Edges[to.Name] = to
}

func (g *Graph) Contains(node string) bool {
	_, ok := g.Nodes[node]
	return ok
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
	g.visited[node.Name] = true

	for _, child := range node.Edges {
		if g.visited[child.Name] {
			continue
		}
		g.topologySort(child)
	}
	g.topoSort = append(g.topoSort, node)
}
