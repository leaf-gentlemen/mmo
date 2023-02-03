package core

import (
	"fmt"
	"sync"
)

type AoiCore struct {
	gID     int
	minX    int
	maxX    int
	minY    int
	maxY    int
	players map[int]bool
	rLock   sync.RWMutex
}

func NewAoiCore(gID, minX, maxX, minY, maxY int) *AoiCore {
	return &AoiCore{
		gID:     gID,
		minX:    minX,
		maxX:    maxX,
		minY:    minY,
		maxY:    maxY,
		players: make(map[int]bool),
	}
}

func (a *AoiCore) Add(uid int) {
	defer a.rLock.Unlock()
	a.rLock.Lock()
	a.players[uid] = true
}

func (a *AoiCore) Remove(uid int) {
	defer a.rLock.Unlock()
	a.rLock.Unlock()
	delete(a.players, uid)
}

func (a *AoiCore) GetPlayerIds() []int {
	defer a.rLock.RUnlock()
	a.rLock.RLock()
	playerIds := make([]int, 0)
	for uid := range a.players {
		playerIds = append(playerIds, uid)
	}
	return playerIds
}

func (a *AoiCore) Sting() string {
	return fmt.Sprint("gid:%d  mixX:%d maxX:%d  minY:%d maxY:%d uids:%v",
		a.gID, a.minX, a.maxX, a.minY, a.maxY, a.players,
	)
}
