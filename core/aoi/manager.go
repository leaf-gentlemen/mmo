package aoi

import "fmt"

// Manager
// @Description: 格子管理中心（地图）
type Manager struct {
	MinX  int           // 整个地图 x 轴最小值
	MaxX  int           // 整个地图 x 轴最大值
	MinY  int           // 整个地图 y 轴最小值
	MaxY  int           // 整个地图 y 抽最大值
	CntX  int           // x 抽方向格子数量
	CntY  int           // y 抽方向格子数量
	grids map[int]*Grid // 格子
}

/*
*  列入如场景大小 250 * 250
*  x抽格子数 cntX = 5
*  y抽格子数 cntY = 5
*  gid = idy*cntY + idx
*  idy = gid / cntY
*  idx = gid % cntX
 */

func NewManager(minX, maxX, minY, maxY, cntX, cntY int) *Manager {
	m := &Manager{
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		CntX:  cntX,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}
	m.loadGrid()
	return m
}

// loadGrid
//
//	@Description: 初始化格子
//	@receiver m
func (m *Manager) loadGrid() {
	for idx := 0; idx < m.CntX; idx++ {
		for idy := 0; idy < m.CntY; idy++ {
			gid := idy*m.CntX + idx                    // 格子 id 算法
			gridMinX := m.MinX + m.getWidth()*idx      // 格子最小 x
			gridMaxX := m.MinX + m.getWidth()*(idx+1)  // 格子最大 x
			gridMinY := m.MinY + m.getLength()*idy     // 格子最小 y
			gridMaxY := m.MinY + m.getLength()*(idy+1) // 格子大 y
			m.grids[gid] = NewGrid(gid, gridMinX, gridMaxX, gridMinY, gridMaxY)
		}
	}
}

// getSurroundGridsByGid
//
//	@Description: 通过 gid 等到周边的 grid
//	@receiver m
//	@param gid
//	@return *Grid
func (m *Manager) getSurroundGridsByGid(gid int) (grids []*Grid) {
	grid, ok := m.grids[gid]
	if !ok {
		return nil
	}

	grids = append(grids, grid)
	// 获取 x 轴编号
	idx := gid % m.CntX
	// 左边不是边界值则取做左侧格子属性
	if grid, ok = m.grids[gid-1]; ok && idx > 0 {
		grids = append(grids, grid)
	}
	// 右边非边界值则取右边的格子属性
	if grid, ok = m.grids[gid+1]; ok && idx < m.CntX-1 {
		grids = append(grids, grid)
	}

	for _, g := range grids {
		idy := g.GID / m.CntX
		// 上边非边界值则取上边的格子属性
		if grid, ok = m.grids[g.GID-m.CntX]; ok && idy > 0 {
			grids = append(grids, grid)
		}
		// 下边非边界值则取取下边的格子属性
		if grid, ok = m.grids[g.GID+m.CntX]; ok && idy < m.CntY-1 {
			grids = append(grids, grid)
		}
	}
	return
}

// String
//
//	@Description: 打印信息
//	@receiver m
//	@return string
func (m *Manager) String() string {
	s := fmt.Sprintf("manager minX:%d maxX:%d minY:%d maxX:%d cntX:%d cntY:%d",
		m.MinX, m.MinX, m.MinY, m.MaxY, m.CntX, m.CntY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// GetGIDByPos
//
//	@Description: 通过坐标获取 gid
//	@receiver m
//	@param x
//	@param y
//	@return int
func (m *Manager) GetGIDByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.getWidth()
	idy := (int(y) - m.MinY) / m.getLength()
	return m.CntX*idy + idx
}

// GetPIDsByPos
//
//	@Description: 通过坐标获取周边的玩家 id
//	@receiver m
//	@param x
//	@param y
//	@return ids
func (m *Manager) GetPIDsByPos(x, y float32) (ids []int) {
	gid := m.GetGIDByPos(x, y)
	grids := m.getSurroundGridsByGid(gid)
	for _, grid := range grids {
		ids = append(ids, grid.GetPlayerIds()...)
	}
	return
}

// GetPIdsByGid
//
//	@Description: 通过 gid 获取格子里的玩家 id
//	@receiver m
//	@param gid
//	@return PIds
func (m *Manager) GetPIdsByGid(gid int) (pIds []int) {
	if grid, ok := m.grids[gid]; ok {
		return grid.GetPlayerIds()
	}
	return nil
}

// RemovePidFromGrid
//
//	@Description:
//	@receiver m
//	@param pid
//	@param gid
func (m *Manager) RemovePidFromGrid(pid, gid int) {
	if grid, ok := m.grids[gid]; ok {
		grid.Remove(pid)
	}
}

// AddPidToGrid
//
//	@Description:
//	@receiver m
//	@param pid
//	@param gid
func (m *Manager) AddPidToGrid(pid, gid int) {
	if grid, ok := m.grids[gid]; ok {
		grid.Add(pid)
	}
}

// RemoveGidByPos
//
//	@Description: 通过坐标移除格子的玩家 id
//	@receiver m
//	@param pid
//	@param x
//	@param y
func (m *Manager) RemoveGidByPos(pid int, x, y float32) {
	gid := m.GetGIDByPos(x, y)
	if grid, ok := m.grids[gid]; ok {
		grid.Remove(pid)
	}
}

// AddPidToPos
//
//	@Description: 通过坐标添加一个玩家 id 到格子里
//	@receiver m
//	@param pid
//	@param x
//	@param y
func (m *Manager) AddPidToPos(pid int, x, y float32) {
	gid := m.GetGIDByPos(x, y)
	if grid, ok := m.grids[gid]; ok {
		grid.Add(pid)
	}
}

// getWidth
//
//	@Description: 获取格子宽
//	@receiver m
//	@return int
func (m *Manager) getWidth() int {
	return (m.MaxX - m.MinX) / m.CntX
}

// getLength
//
//	@Description: 获取格子长
//	@receiver m
//	@return int
func (m *Manager) getLength() int {
	return (m.MaxY - m.MinY) / m.CntY
}
