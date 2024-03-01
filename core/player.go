package core

import (
	"github.com/gogo/protobuf/proto"
	"github.com/leaf-gentlemen/mmo/constants"
	"github.com/leaf-gentlemen/mmo/protos/pubproto"
	"github.com/leaf-gentlemen/zinx/utils"
	"github.com/leaf-gentlemen/zinx/ziface"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math/rand"
	"sync"
)

type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 当前玩家的连接
	X    float32            // 平面x坐标
	Y    float32            // 高度
	Z    float32            // 平面y坐标 (注意不是Y)
	V    float32            // 旋转0-360度
}

var PidGen int32 = 1  // 用来生成玩家ID的计数器
var IdLock sync.Mutex // 保护PidGen的互斥机制

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // 随机在160坐标点 基于X轴偏移若干坐标
		Y:    0,                            // 高度为0
		Z:    float32(134 + rand.Intn(17)), // 随机在134坐标点 基于Y轴偏移若干坐标
		V:    0,                            // 角度为0，尚未实现
	}

	return p
}

// SendMsg
//
//	@Description: 玩家发送消息
//	@receiver p
//	@param msgID 消息ID 可以用路由替换
//	@param data
//	@return error
func (p *Player) SendMsg(msgID uint32, data proto.Message) error {
	buf, err := proto.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := p.Conn.SendMsg(msgID, buf); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// BroadcastStartPosition
//
//	@Description: 广播起始坐标
//	@receiver p
func (p *Player) BroadcastStartPosition() {
	msg := &pubproto.BroadCast{
		Pid: p.Pid,
		Tp:  constants.BroadcastTypeStartPosition,
		Data: &pubproto.BroadCast_P{
			P: &pubproto.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	if err := p.SendMsg(constants.MsgBroadcastPosition, msg); err != nil {
		utils.Logger.Error("BroadcastStartPosition error: ", zap.Error(err))
	}
}

// SynUid
//
//	@Description: 同步客户端坐标
//	@receiver p
func (p *Player) SynUid() {
	msg := &pubproto.SyncPid{
		Pid: p.Pid,
	}

	if err := p.SendMsg(constants.MsgSynUID, msg); err != nil {
		utils.Logger.Error("SynUid error: ", zap.Error(err))
	}
}
