package aoi

import (
	"fmt"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager(0, 250, 0, 250, 5, 5)
	fmt.Println(manager)
}

func TestAOIManagerSuroundGridsByGid(t *testing.T) {
	manager := NewManager(0, 250, 0, 250, 5, 5)
	for k, _ := range manager.grids {
		//得到当前格子周边的九宫格
		grids := manager.GetSurroundGridsByGid(k)
		//得到九宫格所有的IDs
		fmt.Println("gid : ", k, " grids len = ", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Printf("grid ID: %d, surrounding grid IDs are %v\n", k, gIDs)
	}
}
