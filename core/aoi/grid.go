package aoi

import (
	"fmt"
	"sync"
)

// Grid
// @Description: 格子属性和操作
type Grid struct {
	GID     int
	MinX    int
	MaxX    int
	MinY    int
	MaxY    int
	players map[int]bool
	rLock   sync.RWMutex
}

func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:     gID,
		MinX:    minX,
		MaxX:    maxX,
		MinY:    minY,
		MaxY:    maxY,
		players: make(map[int]bool),
	}
}

func (a *Grid) Add(uid int) {
	defer a.rLock.Unlock()
	a.rLock.Lock()
	a.players[uid] = true
}

func (a *Grid) Remove(uid int) {
	defer a.rLock.Unlock()
	a.rLock.Unlock()
	delete(a.players, uid)
}

func (a *Grid) GetPlayerIds() []int {
	defer a.rLock.RUnlock()
	a.rLock.RLock()
	playerIds := make([]int, 0)
	for uid := range a.players {
		playerIds = append(playerIds, uid)
	}
	return playerIds
}

func (a *Grid) Sting() string {
	return fmt.Sprintf("gid:%d  mixX:%d maxX:%d  minY:%d maxY:%d uids:%v",
		a.GID, a.MinX, a.MaxX, a.MinY, a.MaxY, a.players,
	)
}
