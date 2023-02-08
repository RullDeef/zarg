package floormaze

import (
	"log"
	"testing"
)

func TestGeneration(t *testing.T) {
	maze := GenFloorMaze("test maze")
	maze.NextRoom()
	log.Printf("%+v", maze)
}
