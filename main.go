package main

import (
	"fmt"

	"github.com/leaf-gentlemen/mmo/core"
	"github.com/leaf-gentlemen/zinx/utils"
	"github.com/leaf-gentlemen/zinx/ziface"
	"github.com/leaf-gentlemen/zinx/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个玩家
	player := core.NewPlayer(conn)
	// 同步当前的PlayerID给客户端， 走MsgID:1 消息
	player.SynUid()
	// 同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadcastStartPosition()

	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}

func main() {
	utils.InitConf("") // 初始化配置
	srv := znet.NewServe("")
	srv.SetOnConnStart(OnConnectionAdd)
	srv.Serve()
}
