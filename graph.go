//Este codigo es una modificacion del archivo graph.go, proveniente del paquete
//"github.com/twmb/algoimpl/go/graph" creado por twmb.
//Solo se modifico la definicion de nodos y lados para contener el costo y
//beneficio, propiedades necesarias para la realizacion de este proyecto

// Implements an adjacency list graph as a slice of generic nodes
// and includes some useful graph functions.
package main

import (
	"errors"
	"fmt"
)

const (
	dequeued = ^(1<<31 - 1)
	unseen   = 0
	seen     = 1
)

// Graph is an adjacency slice representation of a graph. Can be directed or undirected.
type Graph struct {
	nodes []*node
}

type node struct {
	edges         []edge
	reversedEdges []edge
	index         int
	state         int   // used for metadata
	data          int   // also used for metadata
	parent        *node // also used for metadata
	container     Node  // who holds me
}

// Node connects to a backing node on the graph. It can safely be used in maps.
type Node struct {
	// In an effort to prevent access to the actual graph
	// and so that the Node type can be used in a map while
	// the graph changes metadata, the Node type encapsulates
	// a pointer to the actual node data.
	node *node
	// Value can be used to store information on the caller side.
	// Its use is optional. See the Topological Sort example for
	// a reason on why to use this pointer.
	// The reason it is a pointer is so that graph function calls
	// can test for equality on Nodes. The pointer wont change,
	// the value it points to will. If the pointer is explicitly changed,
	// graph functions that use Nodes will cease to work.
	Value *interface{}
}

type edge struct {
	cost    int
	benefit int
	end     *node
}

// An Edge connects two Nodes in a graph. To modify Weight, use
// the MakeEdgeWeight function. Any local modifications will
// not be seen in the graph.
type Edge struct {
	Cost    int
	Benefit int
	Start   Node
	End     Node
}

// NewGraph creates and returns an empty graph.
// This function returns an undirected graph by default.
func NewGraph() *Graph {
	g := &Graph{}
	return g
}

// MakeNode creates a node, adds it to the graph and returns the new node.
func (g *Graph) MakeNode() Node {
	newNode := &node{index: len(g.nodes)}
	newNode.container = Node{node: newNode, Value: new(interface{})}
	g.nodes = append(g.nodes, newNode)
	return newNode.container
}

// RemoveNode removes a node from the graph and all edges connected to it.
// This function nils points in the Node structure. If 'remove' is used in
// a map, you must delete the map index first.
func (g *Graph) RemoveNode(remove *Node) {
	if remove.node == nil {
		return
	}
	// O(V)
	nodeExists := false
	// remove all edges that connect from a different node to this one
	for _, node := range g.nodes {
		if node == remove.node {
			nodeExists = true
			continue
		}

		// O(E)
		swapIndex := -1 // index that the edge-to-remove is at
		for i := range node.edges {
			if node.edges[i].end == remove.node {
				swapIndex = i
			}
		}
		if swapIndex > -1 {
			swapNRemoveEdge(swapIndex, &node.edges)
		}

		// deal with possible reversed edge
		swapIndex = -1
		for i := range node.reversedEdges {
			if node.reversedEdges[i].end == remove.node {
				swapIndex = i
			}
		}
		if swapIndex > -1 {
			swapNRemoveEdge(swapIndex, &node.reversedEdges)
		}

		if node.index > remove.node.index {
			node.index--
		}
	}
	if nodeExists {
		copy(g.nodes[remove.node.index:], g.nodes[remove.node.index+1:])
		g.nodes = g.nodes[:len(g.nodes)-1]
	}
	remove.node.parent = nil
	remove.node = nil
}

// MakeEdge creates  an edge in the graph with a corresponding cost.
// It returns an error if either of the nodes do not belong in the graph.
//
// Calling MakeEdgeWeight multiple times on the same nodes will not create multiple edges;
// this function will update the weight on the node to the new value.
func (g *Graph) MakeEdge(from, to Node, cost, benefit int) error {
	if from.node == nil || from.node.index >= len(g.nodes) || g.nodes[from.node.index] != from.node {
		return errors.New("First node in MakeEdge call does not belong to this graph")
	}
	if to.node == nil || to.node.index >= len(g.nodes) || g.nodes[to.node.index] != to.node {
		return errors.New("Second node in MakeEdge call does not belong to this graph")
	}

	for i := range from.node.edges { // check if edge already exists
		if from.node.edges[i].end == to.node {
			from.node.edges[i].cost = cost
			from.node.edges[i].benefit = benefit

			// If the graph is undirected, fix the to node's cost as well
			if to != from {
				for j := range to.node.edges {
					if to.node.edges[j].end == from.node {
						to.node.edges[j].cost = cost
						to.node.edges[i].benefit = benefit
					}
				}
			}
			return nil
		}
	}
	newEdge := edge{cost: cost, benefit: benefit, end: to.node}
	from.node.edges = append(from.node.edges, newEdge)
	reversedEdge := edge{cost: cost, benefit: benefit, end: from.node} // cost for undirected graph only
	if to != from {
		to.node.edges = append(to.node.edges, reversedEdge)
	}
	return nil
}

// RemoveEdge removes edges starting at the from node and ending at the to node.
// If the graph is undirected, RemoveEdge will remove all edges between the nodes.
func (g *Graph) RemoveEdge(from, to Node) {
	fromEdges := from.node.edges
	toEdges := to.node.edges
	toReversedEdges := to.node.reversedEdges
	for e := range fromEdges { // fix from->to
		if fromEdges[e].end == to.node {
			swapNRemoveEdge(e, &fromEdges)
			from.node.edges = fromEdges
			break
		}
	}
	for e := range toReversedEdges { // fix reversed edges to->from
		if toReversedEdges[e].end == from.node {
			swapNRemoveEdge(e, &toReversedEdges)
			to.node.reversedEdges = toReversedEdges
			break
		}
	}
	if from.node != to.node {
		for e := range toEdges {
			if toEdges[e].end == from.node {
				swapNRemoveEdge(e, &toEdges)
				to.node.edges = toEdges
				break
			}
		}
	}
}

// Neighbors returns a slice of nodes that are reachable from the given node in a graph.
func (g *Graph) Neighbors(n Node) []Node {
	neighbors := make([]Node, 0, len(n.node.edges))
	if g.nodes[n.node.index] == n.node {
		for _, edge := range n.node.edges {
			neighbors = append(neighbors, edge.end.container)
		}
	}
	return neighbors
}

// Swaps an edge to the end of the edges slice and 'removes' it by reslicing.
func swapNRemoveEdge(remove int, edges *[]edge) {
	(*edges)[remove], (*edges)[len(*edges)-1] = (*edges)[len(*edges)-1], (*edges)[remove]
	*edges = (*edges)[:len(*edges)-1]
}

func (g *Graph) bfs(n *node, finishList *[]Node) {
	totalBenefit := 0
	queue := make([]*node, 0, len(n.edges))
	queue = append(queue, n)
	for i := 0; i < len(queue); i++ {
		node := queue[i]
		node.state = seen
		for _, edge := range node.edges {
			if edge.end.state == unseen {
				edge.end.state = seen
				totalBenefit = totalBenefit + edge.benefit - edge.cost
				queue = append(queue, edge.end)
			}
		}
	}
	fmt.Println(totalBenefit)
	*finishList = make([]Node, 0, len(queue))
	for i := range queue {
		*finishList = append(*finishList, queue[i].container)
	}
}

// ConnectedComponents algorithm for an undirected graph
func (g *Graph) ConnectedComponents() [][]Node {
	components := make([][]Node, 0)
	for _, node := range g.nodes {
		if node.state == unseen {
			component := make([]Node, 0)
			g.bfs(node, &component)
			components = append(components, component)
		}
	}
	fmt.Println(len(components))
	return components
}

// ConnectedComponentOfNode returns the connected component of the i
func (g *Graph) ConnectedComponentOfNode(node *node) []Node {
	component := make([]Node, 0)
	g.bfs(node, &component)
	fmt.Println(component)
	return component
}
