package pathfinding

import (
	"container/heap"
	"gorelay/pkg/packets"
	"math"
)

// Node represents a point in the pathfinding grid
type Node struct {
	X, Y     int
	F, G, H  float64
	Walkable bool
	Parent   *Node
	index    int // Used by heap.Interface
}

// NodeUpdate represents a change in walkability for a node
type NodeUpdate struct {
	X, Y     int
	Walkable bool
}

// Pathfinder implements A* pathfinding
type Pathfinder struct {
	width, height int
	nodes         [][]*Node
}

// NewPathfinder creates a new pathfinder instance
func NewPathfinder(width, height int) *Pathfinder {
	p := &Pathfinder{
		width:  width,
		height: height,
		nodes:  make([][]*Node, height),
	}

	// Initialize nodes
	for y := 0; y < height; y++ {
		p.nodes[y] = make([]*Node, width)
		for x := 0; x < width; x++ {
			p.nodes[y][x] = &Node{
				X:        x,
				Y:        y,
				Walkable: true,
			}
		}
	}

	return p
}

// FindPath finds a path between start and end points
func (p *Pathfinder) FindPath(startX, startY, endX, endY int) []*Node {
	// Validate coordinates
	if !p.isValidCoord(startX, startY) || !p.isValidCoord(endX, endY) {
		return nil
	}

	start := p.nodes[startY][startX]
	end := p.nodes[endY][endX]

	// Reset nodes
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			node := p.nodes[y][x]
			node.F = 0
			node.G = 0
			node.H = 0
			node.Parent = nil
		}
	}

	// Initialize open and closed sets
	openSet := &nodeHeap{}
	heap.Init(openSet)
	heap.Push(openSet, start)
	closedSet := make(map[*Node]bool)

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Node)

		if current == end {
			// Path found, reconstruct it
			path := make([]*Node, 0)
			for current != nil {
				path = append([]*Node{current}, path...)
				current = current.Parent
			}
			return path
		}

		closedSet[current] = true

		// Check neighbors
		for _, neighbor := range p.getNeighbors(current) {
			if !neighbor.Walkable || closedSet[neighbor] {
				continue
			}

			gScore := current.G + p.distance(current, neighbor)
			inOpenSet := false
			for _, node := range *openSet {
				if node == neighbor {
					inOpenSet = true
					break
				}
			}

			if !inOpenSet || gScore < neighbor.G {
				neighbor.Parent = current
				neighbor.G = gScore
				neighbor.H = p.heuristic(neighbor, end)
				neighbor.F = neighbor.G + neighbor.H

				if !inOpenSet {
					heap.Push(openSet, neighbor)
				}
			}
		}
	}

	// No path found
	return nil
}

// UpdateWalkableNodes updates the walkability of nodes
func (p *Pathfinder) UpdateWalkableNodes(updates []NodeUpdate) {
	for _, update := range updates {
		if p.isValidCoord(update.X, update.Y) {
			p.nodes[update.Y][update.X].Walkable = update.Walkable
		}
	}
}

// Helper methods

func (p *Pathfinder) isValidCoord(x, y int) bool {
	return x >= 0 && x < p.width && y >= 0 && y < p.height
}

func (p *Pathfinder) getNeighbors(node *Node) []*Node {
	neighbors := make([]*Node, 0, 8)
	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			if x == 0 && y == 0 {
				continue
			}

			newX := node.X + x
			newY := node.Y + y

			if p.isValidCoord(newX, newY) {
				neighbors = append(neighbors, p.nodes[newY][newX])
			}
		}
	}
	return neighbors
}

func (p *Pathfinder) distance(a, b *Node) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func (p *Pathfinder) heuristic(a, b *Node) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Abs(dx) + math.Abs(dy)
}

// nodeHeap implements heap.Interface for Node priority queue
type nodeHeap []*Node

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].F < h[j].F }
func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *nodeHeap) Push(x interface{}) {
	n := len(*h)
	node := x.(*Node)
	node.index = n
	*h = append(*h, node)
}

func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*h = old[0 : n-1]
	return node
}

// WorldToGrid converts world coordinates to grid coordinates
func WorldToGrid(worldX, worldY float32) (int, int) {
	return int(math.Floor(float64(worldX))), int(math.Floor(float64(worldY)))
}

// GridToWorld converts grid coordinates to world coordinates (centered in tile)
func GridToWorld(gridX, gridY int) (float32, float32) {
	return float32(gridX) + 0.5, float32(gridY) + 0.5
}

// FindPathWorld finds a path between world positions and returns world coordinates
func (p *Pathfinder) FindPathWorld(startPos, endPos *packets.WorldPosData) []*packets.WorldPosData {
	// Convert world to grid coordinates
	startX, startY := WorldToGrid(startPos.X, startPos.Y)
	endX, endY := WorldToGrid(endPos.X, endPos.Y)

	// Find path in grid coordinates
	nodePath := p.FindPath(startX, startY, endX, endY)
	if nodePath == nil {
		return nil
	}

	// Convert back to world coordinates
	worldPath := make([]*packets.WorldPosData, len(nodePath))
	for i, node := range nodePath {
		worldX, worldY := GridToWorld(node.X, node.Y)
		worldPath[i] = &packets.WorldPosData{X: worldX, Y: worldY}
	}

	return worldPath
}
