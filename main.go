package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
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
	beginning := time.Now()

	g := NewGraph()
	positiveG := NewGraph()
	nodes := make(map[int]Node, 0)
	pNodes := make(map[int]Node, 0)
	sortedEdges := Edges{}
	sortedPositiveEdges := Edges{}
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
	args := os.Args
	file, _ := os.Open(args[1])
	optimum, _ := strconv.ParseInt(args[3], 0, 0)
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
			newPositiveEdge := Edge{int(cost), int(benefit), pNodes[int(startNode)], pNodes[int(endNode)]}
			// fmt.Println("Lel")
			// fmt.Println(newEdge)
			// fmt.Println(newPositiveEdge)
			sortedEdges = append(sortedEdges, newEdge)
			sortedPositiveEdges = append(sortedPositiveEdges, newPositiveEdge)
		}
		line++
	}

	sort.Sort(sort.Reverse(sortedPositiveEdges))
	// fmt.Println(sortedEdges)
	g.GraphBuilder(sortedEdges)
	positiveG.PositiveGraphBuilder(sortedPositiveEdges)
	fmt.Println("Grafo Original: \n", g)
	fmt.Println("Grafo Nodos Positivos: \n", positiveG)
	// fmt.Println(positiveG)
	positiveG.unseeNodes()
	positiveG.LinkComponents(sortedPositiveEdges)
	fmt.Println("Grafo Nodos Positivos Conectado: \n", positiveG)
	// Get Floyd Warshall for the complete Graph
	minCost, minPath := g.FloydWarshall()
	//fmt.Println("FW matrix: ", minCost)

	// W need to connect Connected Componentes and get oddNodes

	// positiveG.LinkComponents(sortedEdges)

	positiveG.unseeNodes()

	// Get oddNodes
	oddNodes := make([]int, 0) // List of OddNodes
	for index, elem := range pNodes {
		if positiveG.Degree(elem)%2 != 0 {
			oddNodes = append(oddNodes, index)
		}
	}

	// Compute minimum Matching using Munkres Algorithm
	// Munkres, convert matrix to single vector Munkres Algorithm for OddNodes
	size := len(oddNodes)
	m := mk.NewMatrix(size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			m.A[i*size+j] = int64(minCost[oddNodes[i]-1][oddNodes[j]-1])
		}
		m.A[i*size+i] = math.MaxInt32 // Set infinite to the same Vertice
	}

	minMatching := mk.ComputeMunkresMin(m)
	newMinMatching := []mk.RowCol{}
	minMatchMap := make(map[int]map[int]int)
	for _, elem := range minMatching {
		minMatchMap[oddNodes[elem.Start()]] = make(map[int]int)
		for _ = range positiveG.nodes[elem.Start()].edges {
			minMatchMap[oddNodes[elem.Start()]][oddNodes[elem.End()]] = 1
			// fmt.Println(minMatchMap)
		}

	}
	for _, elem := range minMatching {
		// fmt.Println(minMatchMap[oddNodes[elem.End()]][oddNodes[elem.Start()]])
		if minMatchMap[oddNodes[elem.End()]][oddNodes[elem.Start()]] != 0 {
			minMatchMap[oddNodes[elem.Start()]][oddNodes[elem.End()]] = 0
			// fmt.Println(minMatchMap[oddNodes[elem.End()]][oddNodes[elem.Start()]])
		}
	}
	// fmt.Println("MinMatchMap: ", minMatchMap)
	fmt.Println("Pares a conectar: ")
	for _, elem := range minMatching {
		if minMatchMap[oddNodes[elem.Start()]][oddNodes[elem.End()]] != 0 {
			newMinMatching = append(newMinMatching, elem)
			fmt.Print("(", oddNodes[elem.Start()], ",", oddNodes[elem.End()], "),")
		}
	}
	fmt.Println()
	// Insert Path from Munkres algorithm
	for _, elem := range newMinMatching {
		startIndex := oddNodes[elem.Start()]
		start := nodes[startIndex].node
		path := ReconstructPath(minPath, oddNodes[elem.Start()]-1, oddNodes[elem.End()]-1)
		for _, vertice := range path {
			nextIndex := vertice + 1
			next := nodes[nextIndex].node
			for _, edge := range start.edges {
				if edge.end == next {
					fmt.Printf("Agregando Arista: (%d,%d)\n", startIndex, vertice+1)
					positiveG.MakeEdge(pNodes[startIndex], pNodes[nextIndex], edge.cost, edge.benefit)
					break
				}
			}
			start = next
			startIndex = nextIndex

		}
	}
	fmt.Println("Imprimiendo grafo positivo nuevo con todos los nodos pares")
	fmt.Println(positiveG)
	eulerPath, _, value := positiveG.EulerianCycle(pNodes[1])
	fmt.Println(eulerPath)
	salida, err := os.Create(args[2] + "-salida.txt")
	check(err)
	defer salida.Close()
	stringValue := strconv.Itoa(value)
	stringPath := []string{}
	_, err = salida.WriteString(stringValue)
	check(err)
	_, err = salida.WriteString("\n")
	check(err)
	for i := range eulerPath {
		number := eulerPath[len(eulerPath)-i-1]
		text := strconv.Itoa(number)
		stringPath = append(stringPath, text)
	}
	result := strings.Join(stringPath, " ")
	_, err = salida.WriteString(result)
	check(err)
	salida.Sync()

	optimumDeviation := float64(100 * ((int(optimum)) - value) / int(optimum))

	fmt.Println(optimum)

	elapsed := time.Since(beginning)
	fmt.Println(elapsed)
	fmt.Println(optimumDeviation)
}
