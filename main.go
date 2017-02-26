package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	// "github.com/twmb/algoimpl/go/graph"
	"strings"

	mk "./munkres"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Hello World")

	g := NewGraph()
	positiveG := NewGraph()
	nodes := make(map[int]Node, 0)
	pNodes := make(map[int]Node, 0)
	sortedEdges := Edges{}
	// nodes[0] = g.MakeNode()
	// nodes[1] = g.MakeNode()
	// nodes[2] = g.MakeNode()
	// nodes[3] = g.MakeNode()
	// nodes[4] = g.MakeNode()
	// nodes[5] = g.MakeNode()
	// nodes[6] = g.MakeNode()
	// nodes[0] = g.MakeNode()
	// g.MakeEdge(nodes[0], nodes[1], 2, 10)
	// g.MakeEdge(nodes[0], nodes[2], 10, 0)
	// g.MakeEdge(nodes[1], nodes[2], 3, 2)
	// g.MakeEdge(nodes[1], nodes[3], 20, 5)
	// g.MakeEdge(nodes[1], nodes[4], 1, 3)
	// g.MakeEdge(nodes[2], nodes[3], 3, 4)
	// g.MakeEdge(nodes[2], nodes[4], 5, 4)
	// g.MakeEdge(nodes[3], nodes[4], 2, 8)
	// g.MakeEdge(nodes[3], nodes[5], 9, 1)
	// g.MakeEdge(nodes[4], nodes[5], 8, 1)

	file, _ := os.Open("./test")
	lineScanner := bufio.NewScanner(file)
	line := 0
	for lineScanner.Scan() {
		contents := strings.Fields(lineScanner.Text())
		if line == 0 {
			number, _ := strconv.ParseInt(contents[len(contents)-1], 0, 0)
			for i := 1; i < int(number+1); i++ {
				nodes[i] = g.MakeNode()
				*nodes[i].Value = i
				pNodes[i] = positiveG.MakeNode()
				*pNodes[i].Value = i
			}
		}
		if _, err := strconv.Atoi(contents[0]); err == nil {
			startNode, _ := strconv.ParseInt(contents[0], 0, 0)
			endNode, _ := strconv.ParseInt(contents[1], 0, 0)
			cost, _ := strconv.ParseInt(contents[2], 0, 0)
			benefit, _ := strconv.ParseInt(contents[3], 0, 0)
			newEdge := Edge{int(cost), int(benefit), nodes[int(startNode)], nodes[int(endNode)]}
			// fmt.Println(newEdge)
			sortedEdges = append(sortedEdges, newEdge)
		}
		line++
	}
	sort.Sort(sort.Reverse(sortedEdges))
	// fmt.Println(sortedEdges)
	g.GraphBuilder(sortedEdges)
	positiveG.PositiveGraphBuilder(sortedEdges)
	// g.ConnectedComponentsMap()
	g.unseeNodes()
	g.checkIncidence()
	// sort.Reverse(sortedEdges)
	eulerPath, _ := g.EulerianCycle(nodes[1])
	fmt.Println(eulerPath)
	// g.ConnectedComponentOfNode(nodes[1].node)
	// g.ConnectedComponents()
	// fmt.Println(nodes[1])
	// check(err)

	// Get Floyd Warshall for the complete Graph
	minPath := g.FloydWarshall()

	// W need to connect Connected Componentes and get oddNodes

	// Get oddNodes
	oddNodes := make([]int, 0)
	for index, elem := range pNodes {
		if positiveG.Degree(elem)%2 != 0 {
			oddNodes = append(oddNodes, index)
		}
	}

	// Munkres, convert matrix to single vector Munkres Algorithm for OddNodes
	size := len(oddNodes)
	m := mk.NewMatrix(size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			m.A[i*size+j] = int64(minPath[oddNodes[i]-1][oddNodes[j]-1])
		}
		m.A[i*size+i] = math.MaxInt32 // Set infinite to the same Vertice
	}
	fmt.Println("FW matrix: ", minPath)
	fmt.Println(m.A)
	fmt.Println(mk.ComputeMunkresMin(m)) // [{0 3} {1 1} {2 0} {3 2}]

}
