package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/frame/g"
)

const (
	PkgRbpIos     = "com.im.duck.ios"
	PkgRbpAndroid = "com.im.android.rbp"
	PkgRbpPc      = "com.im.rbp.desktop"
)

// 支付、充值相关
const (
	CmdIapGoogle = "iap.google" // 谷歌的iap
	CmdIapApple  = "iap.apple"  //苹果的的iap

	CmdActivityHandle = "activity.handle" //作品

	IapTypeGiap     = "giap"     //google的type
	IapTypeAiap     = "iap"      //apple的type
	IapTypeOfficial = "official" //官方充值
	TopicSLPIap     = "slp.iap"  //iap的topic

	// IapProductTypeRecharge 产品的类型是充值
	IapProductTypeRecharge = "recharge"
	// IapProductTypeTitle 产品的类型爵位续费，具体续费多少有product_extra_val确定，单位为月
	IapProductTypeTitle = "title"

	IapHandleStatusFinished = uint32(1) //iap处理完了，不需要处理
	IapHandleStatusInit     = uint32(0) //ipa没有处理，需要进一步处理
)

// IsDev 表明当前系统是不是dev
// 不要在init初始化函数中使用
var IsDev bool = false
var RunMode string = "prod"

func init() {
	//通过机器来检测是不是alpha，这样来统一所有的配置和部署
	mode := g.Cfg().GetString("server.RunMode")
	alphaHosts := g.Cfg().GetStrings("server.AlphaHosts")
	if alphaHosts != nil && mode == "prod" {
		host, err := os.Hostname()
		if err == nil {
			host = strings.ToLower(host)
			for _, alpha := range alphaHosts {
				if host == strings.ToLower(alpha) {
					//是alpha服务器
					mode = "alpha"
				}
			}
		}
	}

	IsDev = mode == "dev" || len(mode) == 0
	RunMode = mode

	fmt.Println("server run with", RunMode)
}

// AppID 不同独立APP的唯一标识
type AppID uint8

// 定义不同APP
const (
	AppUnknown       AppID = 0  //未备案
	AppRBP           AppID = 88 // 彩虹星球
	AppChongYa       AppID = 2  //冲鸭
	AppXiongShou     AppID = 3  //凶手
	AppKaiXin        AppID = 4  //开心玩
	AppOverseaBanBan AppID = 5  //PT
	AppRush          AppID = 6  //冲鸭
	AppMax           AppID = 7  //
	AppMinGrp        AppID = 8  //星球小程序动态过滤虚拟app_id
	AppChongYaDummy  AppID = 12 //冲鸭马甲包，只为ios过审使用
	AppYinLang       AppID = 67 // 音浪派对

)

// IntToAppID 把数字类型转成 AppID 类型
func IntToAppID(value int) AppID {
	if value == int(AppChongYaDummy) {
		return AppChongYaDummy
	}
	if value >= 1 && value <= 100 {
		return AppID(value)
	}
	return AppUnknown
}

var packageToAppID map[string]AppID = map[string]AppID{
	"com.im.android.rbp": AppRBP, //彩虹星球安卓
	"com.im.duck.ios":    AppRBP, //彩虹星球安卓
}

// 隔离app显示的配置，不在配置里的不隔离
var isolateAppGroupMap = map[AppID][]AppID{
	AppChongYa:       {AppChongYa},
	AppRush:          {AppRush},
	AppXiongShou:     {AppXiongShou, AppKaiXin},
	AppKaiXin:        {AppXiongShou, AppKaiXin},
	AppMax:           {AppMax},
	AppMinGrp:        {AppMinGrp},
	AppOverseaBanBan: {AppOverseaBanBan},
	AppRBP:           {AppRBP},
	AppYinLang:       {AppYinLang},
}

// 是否是小程序
var miniPackageMap = map[string]AppID{
	"im.banban.mini.weapp":  AppRBP,
	"im.banban4.mini.weapp": AppRBP,
	"im.banban5.mini.weapp": AppRBP,
	"im.banban2.mini.qq":    AppRBP,
}

// GetAppID 把客户端的包名转成对应的AppID
func GetAppID(packageName string) AppID {
	if appID, ok := packageToAppID[packageName]; ok {
		return appID
	}
	return AppUnknown
}

// 根据当前app的app_id和检测目标app_id，判断是否需要做隔离处理
func NeedAppIsolate(curAppId, targetAppId AppID) bool {
	allowAppIds := isolateAppGroupMap[curAppId]
	if len(allowAppIds) == 0 {
		return false
	}
	for _, v := range allowAppIds {
		if v == targetAppId {
			return false
		}
	}
	return true
}

// 根据当前appid，获取app对应分组的app_id列表
func GetIsolateGroupAppIds(appId AppID) []AppID {
	return isolateAppGroupMap[appId]
}

// IsTableAppID 判断是否桌游app_id
func IsTableAppID(appId uint32) bool {
	return IsKiller(appId) || appId == uint32(AppKaiXin)
}

// IsKiller 判断是否桌游app_id
func IsKiller(appId uint32) bool {
	return appId == uint32(AppXiongShou)
}

// IsMate 判断是否皮队友app_id
func IsMate(appId uint32) bool {
	return appId == uint32(AppChongYa)
}

// IsBanBan 判断是否嗨歌app_id
func IsBanBan(appId uint32) bool {
	return appId == uint32(AppRBP)
}

// IsHiSong 判断是否嗨歌app_id
func IsHiSong(packageName string) bool {
	dict := map[string]AppID{
		"com.havefun.android": AppRBP,
		"com.draw.guess.ios":  AppRBP,
	}
	_, ok := dict[packageName]
	return ok
}

// 先局部改一下，后面出需求确认好业务搏击范围，再统一替换  IsHiSong 为 IsSlp
func IsSlp(packageName string) bool {
	dict := map[string]AppID{
		"com.yhl.sleepless.android": AppRBP,
		"sg.ola.party.alo":          AppRBP,
	}
	_, ok := dict[packageName]
	return ok
}

// IsMini 判断是否嗨歌小程序
func IsMini(packageName string) bool {
	_, ok := miniPackageMap[packageName]
	return ok
}

func IsSlpAndroid(packageName string) bool {
	return packageName == "com.yhl.sleepless.android"
}

func IsSlpIos(packageName string) bool {
	return packageName == "sg.ola.party.alo"
}

// IsYlp 判断是否音浪马甲包
func IsYlp(packageName string) bool {
	appId, ok := packageToAppID[packageName]
	if ok && appId == AppYinLang {
		return true
	}
	return false
}

//todo... 给不同app or 小程序 or web 定义不用的query sign
