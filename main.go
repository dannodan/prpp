package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	// "github.com/twmb/algoimpl/go/graph"
	"strings"
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
			newEdge := Edge{int(cost), int(benefit), pNodes[int(startNode)], pNodes[int(endNode)]}
			// fmt.Println("Lel")
			// fmt.Println(newEdge)
			sortedEdges = append(sortedEdges, newEdge)
		}
		line++
	}
	// fmt.Println(nodes)
	sort.Sort(sort.Reverse(sortedEdges))
	// fmt.Println(sortedEdges)
	// g.GraphBuilder(sortedEdges)
	positiveG.GraphBuilder(sortedEdges)
	fmt.Println(positiveG)
	positiveG.ConnectedComponentsMap()
	positiveG.unseeNodes()
	positiveG.LinkComponents(sortedEdges)
	// g.checkIncidence()
	// sort.Reverse(sortedEdges)
	// eulerPath, _ := g.EulerianCycle(nodes[1])
	// fmt.Println(eulerPath)
	// g.ConnectedComponents()
	// fmt.Println(positiveG.nodes[0].container)
	// check(err)

	// fmt.Println(positiveG.ConnectedComponents())

	positiveG.unseeNodes()

	// path := positiveG.GetPath(positiveG.nodes[0].container)

	// fmt.Println(path)
	//
	// // Get Floyd Warshall for the complete Graph
	// minCost, minPath := positiveG.FloydWarshall()
	// fmt.Println("FW matrix: ", minCost)
	//
	// // W need to connect Connected Componentes and get oddNodes
	//
	// // positiveG.LinkComponents(sortedEdges)
	//
	// positiveG.unseeNodes()
	//
	// fmt.Println(positiveG.ConnectedComponents())
	//
	// // Get oddNodes
	// oddNodes := make([]int, 0) // List of OddNodes
	// for index, elem := range pNodes {
	// 	if positiveG.Degree(elem)%2 != 0 {
	// 		oddNodes = append(oddNodes, index)
	// 	}
	// }
	//
	// fmt.Println()
	// fmt.Println("Imprimiendo grafo positivo Original")
	// fmt.Println(positiveG)
	//
	// // Compute minimum Matching using Munkres Algorithm
	// // Munkres, convert matrix to single vector Munkres Algorithm for OddNodes
	// size := len(oddNodes)
	// m := mk.NewMatrix(size)
	// for i := 0; i < size; i++ {
	// 	for j := 0; j < size; j++ {
	// 		m.A[i*size+j] = int64(minCost[oddNodes[i]-1][oddNodes[j]-1])
	// 	}
	// 	m.A[i*size+i] = math.MaxInt32 // Set infinite to the same Vertice
	// }
	//
	// minMatching := mk.ComputeMunkresMin(m)
	// newMinMatching := []mk.RowCol{}
	// fmt.Println(m.A)
	// minMatchMap := make(map[int]map[int]int)
	// for _, elem := range minMatching {
	// 	minMatchMap[oddNodes[elem.Start()]] = make(map[int]int)
	// 	for _ = range positiveG.nodes[elem.Start()].edges {
	// 		minMatchMap[oddNodes[elem.Start()]][oddNodes[elem.End()]] = 1
	// 	}
	// 	fmt.Print("(", oddNodes[elem.Start()], ",", oddNodes[elem.End()], "), ")
	// }
	// for _, elem := range minMatching {
	// 	if minMatchMap[elem.End()][elem.Start()] != 0 {
	// 		minMatchMap[elem.End()][elem.Start()] = 0
	// 	}
	// }
	// for _, elem := range minMatching {
	// 	if minMatchMap[elem.Start()][elem.Start()] != 0 {
	// 		newMinMatching = append(newMinMatching, elem)
	// 	}
	// }
	//
	// fmt.Println()
	// // Insert Path from Munkres algorithm
	// for _, elem := range newMinMatching {
	// 	startIndex := oddNodes[elem.Start()]
	// 	start := nodes[startIndex].node
	// 	fmt.Println(oddNodes[elem.Start()])
	// 	path := ReconstructPath(minPath, oddNodes[elem.Start()]-1, oddNodes[elem.End()]-1)
	// 	for _, vertice := range path {
	// 		nextIndex := vertice + 1
	// 		next := nodes[nextIndex].node
	// 		for _, edge := range start.edges {
	// 			if edge.end == next {
	// 				// fmt.Printf("Agregando Arista: (%d,%d)\n", startIndex, vertice+1)
	// 				// fmt.Println(startIndex-1)
	// 				// fmt.Println(pNodes[startIndex])
	// 				// fmt.Println(nextIndex)
	// 				// fmt.Println(pNodes[nextIndex])
	// 				// positiveG.checkEdges(pNodes[startIndex], pNodes[nextIndex], edge.cost, edge.benefit)
	// 				positiveG.MakeEdge(pNodes[startIndex], pNodes[nextIndex], edge.cost, edge.benefit)
	// 				break
	// 			}
	// 		}
	// 		start = next
	// 		startIndex = nextIndex
	//
	// 	}
	// }
	//
	// // fmt.Println(positiveG.ConnectedComponentOfNode(nodes[1].node))
	// // positiveG.checkIncidence()
	// // fmt.Println()
	// // fmt.Println("Imprimiendo grafo positivo nuevo")
	fmt.Println(positiveG)
	// // eulerPath, _ := positiveG.EulerianCycle(nodes[1])
	// // fmt.Println(eulerPath)
	// // check(err)
}
