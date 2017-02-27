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
	"math"
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
	incidence     int   // used for incidence
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
	state   int
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

type Edges []Edge

func (slice Edges) Len() int {
	return len(slice)
}

func (slice Edges) Less(i, j int) bool {
	return (float64(slice[i].Benefit) / float64(slice[i].Cost)) < (float64(slice[j].Benefit) / float64(slice[j].Cost))
}

func (slice Edges) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// To String Function
func (n Node) String() string {
	return fmt.Sprintf("%d", *n.Value)
}

func (n *node) String() string {
	ed := fmt.Sprintf("%d -> [", *n.container.Value)
	for _, edge := range n.edges {
		ed = ed + fmt.Sprintf("(%d,%d)", *n.container.Value, *edge.end.container.Value)
	}
	ed = ed + "]"
	return ed
}

func (g Graph) String() string {
	nodes := ""
	for _, node := range g.nodes {
		nodes = nodes + node.String() + "\n"
	}
	return nodes
}

// NewGraph creates and returns an empty graph.
// This function returns an undirected graph by default.
func NewGraph() *Graph {
	g := &Graph{}
	return g
}

// MakeNode creates a node, adds it to the graph and returns the new node.
func (g *Graph) MakeNode() Node {
	newNode := &node{index: len(g.nodes), incidence: 0}
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
	// fmt.Println(from.node)
	if from.node == nil || from.node.index >= len(g.nodes) || g.nodes[from.node.index] != from.node {
		return errors.New("First node in MakeEdge call does not belong to this graph")
	}
	if to.node == nil || to.node.index >= len(g.nodes) || g.nodes[to.node.index] != to.node {
		// fmt.Println("Error MakeEdge")
		return errors.New("Second node in MakeEdge call does not belong to this graph")
	}

	// for i := range from.node.edges { // check if edge already exists
	// 	if from.node.edges[i].end == to.node {
	// 		from.node.edges[i].cost = cost
	// 		from.node.edges[i].benefit = benefit
	//
	// 		// If the graph is undirected, fix the to node's cost as well
	// 		if to != from {
	// 			for j := range to.node.edges {
	// 				if to.node.edges[j].end == from.node {
	// 					to.node.edges[j].cost = cost
	// 					to.node.edges[i].benefit = benefit
	// 				}
	// 			}
	// 		}
	// 		return nil
	// 	}
	// }
	newEdge := edge{cost: cost, benefit: benefit, end: to.node}
	from.node.edges = append(from.node.edges, newEdge)
	reversedEdge := edge{cost: cost, benefit: benefit, end: from.node} // cost for undirected graph only
	if to != from {
		to.node.edges = append(to.node.edges, reversedEdge)
	}
	// fmt.Println("Lado Creado")
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
			fmt.Println(fromEdges[e])
			swapNRemoveEdge(e, &fromEdges)
			from.node.edges = fromEdges
			// fmt.Println("erased")
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
	fmt.Println(edges)
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
	// fmt.Println(totalBenefit)
	*finishList = make([]Node, 0, len(queue))
	for i := range queue {
		*finishList = append(*finishList, queue[i].container)
	}
}

func (g *Graph) bfsMap(n *node, finishList map[int]Node) map[int]Node {
	queue := make([]*node, 0, len(n.edges))
	queue = append(queue, n)
	for i := 0; i < len(queue); i++ {
		node := queue[i]
		node.state = seen
		for _, edge := range node.edges {
			if edge.end.state == unseen {
				edge.end.state = seen
				queue = append(queue, edge.end)
			}
		}
	}
	finishList = make(map[int]Node, len(queue))
	for i := range queue {
		// fmt.Println(queue[i].index)
		finishList[queue[i].index] = queue[i].container
	}
	return finishList
}

func (g *Graph) ConnectedComponentsMap() []map[int]Node {
	componentMap := make([]map[int]Node, 0)
	for _, node := range g.nodes {
		if node.state == unseen {
			component := make(map[int]Node, 0)
			component = g.bfsMap(node, component)
			componentMap = append(componentMap, component)
		}
	}
	// fmt.Println(componentMap)
	return componentMap
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

func (g *Graph) LinkComponents(edges Edges) {
	// linkedComponents := make([]map[int]Node, 0)
	linkedComponents := g.ConnectedComponentsMap()
	g.unseeNodes()
	// fmt.Println(g)
	//fmt.Println(linkedComponents)
	for _, edge := range edges {
		// fmt.Println(edge)
		for _, component := range linkedComponents {
			first := component[edge.Start.node.index]
			second := component[edge.End.node.index]
			if ((first != Node{}) && (second != Node{})) || ((first == Node{}) && (second == Node{})) {
				// fmt.Println("Do Nothing")
			} else {
				// fmt.Println(first.node.index)
				// fmt.Println(second.node.index)
				g.MakeEdge(edge.Start, edge.End, edge.Cost, edge.Benefit)
				// fmt.Println("Nodo Agregado")
				break
			}
		}
		linkedComponents = g.ConnectedComponentsMap()
		g.unseeNodes()
		//fmt.Println(g)
	}
}

func (g *Graph) GraphBuilder(edges Edges) {
	//fmt.Println(edges)
	for _, edge := range edges {
		g.MakeEdge(edge.Start, edge.End, edge.Cost, edge.Benefit)
		edge.Start.node.incidence = edge.Start.node.incidence + 1
		edge.End.node.incidence = edge.End.node.incidence + 1
	}
}

func (g *Graph) PositiveGraphBuilder(edges Edges) {
	//fmt.Println(edges)
	for _, edge := range edges {
		if edge.Benefit-edge.Cost >= 0 {
			g.MakeEdge(edge.Start, edge.End, edge.Cost, edge.Benefit)
			edge.Start.node.incidence = edge.Start.node.incidence + 1
			edge.End.node.incidence = edge.End.node.incidence + 1
		}
	}
}

func (g *Graph) unseeNodes() {
	for _, node := range g.nodes {
		node.state = unseen
	}
}

func (g *Graph) checkIncidence() {
	totalNodes := 0
	for _, node := range g.nodes {
		if node.incidence%2 != 0 {
			totalNodes++
		}
	}
	// fmt.Println(totalNodes)
}

// func (g *Graph) GetPath(fromNode Node) []int {
// 	path := []int{}
// 	edgeList := make(map[Node]map[Node]edge, 0)
// 	for _, node := range g.nodes {
// 		edgeList[node.container] = make(map[Node]edge, 0)
// 		for _, edge := range node.edges {
// 			edgeList[node.container][edge.end.container] = edge
// 		}
// 	}
// 	chosenEdge := 0
// 	costIndex := float64(0)
// 	for index, edge := range fromNode.node.edges {
// 		test := float64(edgeList[fromNode][edge.end.container].benefit) / float64(edgeList[fromNode][edge.end.container].cost)
// 		if (test > costIndex) && edge.state <= fromNode.node.edges[chosenEdge].state {
// 			fmt.Println(edge.state)
// 			costIndex = test
// 			chosenEdge = index
// 			// fmt.Println("test")
// 			// fmt.Println(chosenEdge)
// 		}
// 	}
// 	fmt.Println("test")
// 	edgeList[fromNode][fromNode.node.edges[chosenEdge].end.container].state = edgeList[fromNode][fromNode.node.edges[chosenEdge].end.container].state + 1
// 	// edgeList[fromNode][fromNode.node.edges[chosenEdge].end.container]
// 	fmt.Println(fromNode.node.edges[chosenEdge].state)
// 	if fromNode.node.edges[chosenEdge].end == g.nodes[0] {
// 		path = append(path, 1)
// 		fmt.Println(path)
// 	} else {
// 		path = g.GetPath(fromNode.node.edges[chosenEdge].end.container)
// 		path = append(path, fromNode.node.edges[chosenEdge].end.index)
// 		fmt.Println(path)
// 	}
// 	return path
// }

func (g *Graph) EulerianCycle(start Node) (tour []int, success bool) {
	// For an Eulerian cirtuit all the vertices has to have a even degree
	// if start.node.incidence < 2 {
	// 	fmt.Println(start.node.edges[0].end.container)
	// 	g.MakeEdge(start, start.node.edges[0].end.container, start.node.edges[0].cost, 0)
	// 	start.node.incidence++
	// }
	unvisitedEdges := make(map[Node]map[Node]int, 0)
	for _, node := range g.nodes {
		if len(node.edges)%2 != 0 {
			return nil, false
		}
		unvisitedEdges[node.container] = make(map[Node]int, 0)
		for _, edge := range node.edges {
			unvisitedEdges[node.container][edge.end.container] = edge.benefit - edge.cost
		}
	}
	fmt.Println(unvisitedEdges)
	// Hierholzer's algorithm
	var currentNode, nextNode Node
	//
	valueStack := []int{}
	value := 0
	tour = []int{}
	stack := []Node{start}
	for len(stack) > 0 {
		currentNode = stack[len(stack)-1]
		// Get an arbitrary edge from the current vertex
		// 	edgesSeen := 0
		// fmt.Println(unvisitedEdges[currentNode])
		// fmt.Println(len(unvisitedEdges[currentNode]))
		if len(unvisitedEdges[currentNode]) > 0 {
			for nextNode = range unvisitedEdges[currentNode] {
				break
			}
			// fmt.Println(unvisitedEdges[currentNode][nextNode])
			valueStack = append(valueStack, unvisitedEdges[currentNode][nextNode])
			delete(unvisitedEdges[currentNode], nextNode)
			delete(unvisitedEdges[nextNode], currentNode)
			stack = append(stack, nextNode)
			fmt.Println(valueStack[len(valueStack)-1])
		} else {
			// fmt.Println(len(valueStack))
			tour = append(tour, stack[len(stack)-1].node.index+1)
			// fmt.Println(value)
			stack = stack[:len(stack)-1]
		}
	}
	for index := range valueStack {
		value = value + valueStack[index]
		// valueStack = valueStack[:len(stack)-1]
	}
	fmt.Println(value)
	return tour, true
}

func (g *Graph) Degree(n Node) int {
	return len(n.node.edges)
}

func (g *Graph) FloydWarshall() (mincost, minpath [][]int) {
	path := make([][]int, len(g.nodes))
	next := make([][]int, len(g.nodes))
	// Build Distance Matrix
	lenNodes := len(g.nodes)
	for i := 0; i < lenNodes; i++ {
		path[i] = make([]int, lenNodes)
		next[i] = make([]int, lenNodes)
		for j := 0; j < lenNodes; j++ {
			path[i][j] = math.MaxInt32
			next[i][j] = -1
		}
		path[i][i] = 0
		for _, edge := range g.nodes[i].edges {
			path[i][edge.end.index] = edge.cost
			next[i][edge.end.index] = edge.end.index
		}
	}
	// Floyd Warshall Algorithm
	for k := 0; k < lenNodes; k++ {
		for i := 0; i < lenNodes; i++ {
			for j := 0; j < lenNodes; j++ {
				dt := path[i][k] + path[k][j]
				if path[i][j] > dt {
					path[i][j] = dt
					next[i][j] = next[i][k]
				}
			}
		}
	}
	return path, next
}

func ReconstructPath(next [][]int, u, v int) []int {
	if next[u][v] == -1 {
		return nil
	}
	path := make([]int, 0)
	for u != v {
		u = next[u][v]
		path = append(path, u)
	}
	return path
}
