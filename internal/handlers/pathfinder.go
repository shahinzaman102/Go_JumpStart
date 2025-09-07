package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shahinzaman102/Go_JumpStart/internal/pathfinder"
)

// PathfinderResponse represents the result of a pathfinding algorithm.
type PathfinderResponse struct {
	Algorithm     string  `json:"algorithm"`
	PathLength    int     `json:"path_length"`
	ExecutionTime float64 `json:"execution_time_ms"`
}

// Pathfinder runs pathfinding algorithms (Brute Force and BFS)
// on a sample grid and returns their performance as JSON.
// Query parameters: start=x,y and end=x,y
func Pathfinder(w http.ResponseWriter, r *http.Request) {
	startParam := r.URL.Query().Get("start")
	endParam := r.URL.Query().Get("end")

	if startParam == "" || endParam == "" {
		http.Error(w, "Missing required params: start, end", http.StatusBadRequest)
		return
	}

	startCoords := strings.Split(startParam, ",")
	endCoords := strings.Split(endParam, ",")

	if len(startCoords) != 2 || len(endCoords) != 2 {
		http.Error(w, "Invalid coordinates format. Use start=x,y end=x,y", http.StatusBadRequest)
		return
	}

	startX, err := strconv.Atoi(startCoords[0])
	if err != nil {
		http.Error(w, "Invalid start X coordinate", http.StatusBadRequest)
		return
	}
	startY, err := strconv.Atoi(startCoords[1])
	if err != nil {
		http.Error(w, "Invalid start Y coordinate", http.StatusBadRequest)
		return
	}

	destX, err := strconv.Atoi(endCoords[0])
	if err != nil {
		http.Error(w, "Invalid destination X coordinate", http.StatusBadRequest)
		return
	}
	destY, err := strconv.Atoi(endCoords[1])
	if err != nil {
		http.Error(w, "Invalid destination Y coordinate", http.StatusBadRequest)
		return
	}

	// Sample grid (0 = free cell, 1 = obstacle)
	grid := [][]int{
		{0, 0, 1, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 1, 0},
		{0, 0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 1, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 1, 0, 0},
	}

	rows := len(grid)
	cols := len(grid[0])

	// Validate bounds
	if startX < 0 || startY < 0 || startX >= rows || startY >= cols {
		http.Error(w, "Invalid start coordinates: out of bounds", http.StatusBadRequest)
		return
	}
	if destX < 0 || destY < 0 || destX >= rows || destY >= cols {
		http.Error(w, "Invalid destination coordinates: out of bounds", http.StatusBadRequest)
		return
	}

	// Validate obstacles
	if grid[startX][startY] == 1 {
		http.Error(w, "Invalid start coordinates: cannot be on an obstacle", http.StatusBadRequest)
		return
	}
	if grid[destX][destY] == 1 {
		http.Error(w, "Invalid destination coordinates: cannot be on an obstacle", http.StatusBadRequest)
		return
	}

	// const runs = 1
	const runs = 1000
	results := []PathfinderResponse{}

	// Brute Force
	var brutePathLen int
	startTime := time.Now()
	for i := 0; i < runs; i++ {
		visited := make([][]bool, len(grid))
		for r := range visited {
			visited[r] = make([]bool, len(grid[0]))
		}
		path := pathfinder.ShortestPathBruteForce(grid, startX, startY, destX, destY, visited)
		if path == math.MaxInt32 {
			path = -1
		}
		brutePathLen = path
	}
	bruteDuration := time.Since(startTime).Seconds() * 1000 / runs

	results = append(results, PathfinderResponse{
		Algorithm:     "brute",
		PathLength:    brutePathLen,
		ExecutionTime: bruteDuration,
	})

	// BFS
	var bfsPathLen int
	startTime = time.Now()
	for i := 0; i < runs; i++ {
		bfsPathLen = pathfinder.ShortestPathBFS(grid, startX, startY, destX, destY)
	}
	bfsDuration := time.Since(startTime).Seconds() * 1000 / runs

	results = append(results, PathfinderResponse{
		Algorithm:     "bfs",
		PathLength:    bfsPathLen,
		ExecutionTime: bfsDuration,
	})

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
