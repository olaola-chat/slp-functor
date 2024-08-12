package model

import (
	"github.com/gogf/gf/frame/g"
)

var RunMode string = "prod"

func init() {
	RunMode = g.Cfg().GetString("server.RunMode")
}

// AppID 不同独立APP的唯一标识
type AppID uint8

// 定义不同APP
const (
	AppRBP AppID = 88 // 彩虹星球
)
