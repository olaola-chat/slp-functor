package model

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/smallnest/rpcx/share"
)

const (
	// ContextKey 上下文变量存储键名，前后端系统共享
	ContextKey = "ContextKey"
	// TrackKey 上下文传递
	TrackKey = "TraceId"
	// TraceingEnabled 当前请求是否开启
	TraceingEnabled = "TraceingEnabled"
)

var languages map[string]bool = make(map[string]bool)
var defaultLang string

// GetSupportLanguages 返回系统当前支持哪些语言
// 不要在init函数里调用
func GetSupportLanguages() []string {
	res := make([]string, len(languages))
	i := 0
	for lang := range languages {
		res[i] = lang
		i = i + 1
	}
	return res
}

func init() {
	//获取系统支持的语言
	ls := g.Cfg().GetStrings("server.Language")
	if len(ls) == 0 {
		panic(gerror.New("empty language in config"))
	}

	defaultLang = strings.ToLower(ls[0])
	for _, lang := range ls {
		languages[strings.ToLower(lang)] = true
	}
	g.Log().Println("init lang languages", ls)
	g.Log().Println("init lang default", defaultLang)
}

// GetContext 根据请求context获取当前信息
func GetContext(ctx context.Context) *Context {
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*Context); ok {
		return localCtx
	}
	return nil
}

func GetOpentracingEnabled(ctx context.Context) bool {
	val := ctx.Value(TraceingEnabled)
	if v, ok := val.(bool); ok {
		return v
	}
	return false
}

// GetOpentracingSpanFromContext 从Context里获取Span
func GetOpentracingSpanFromContext(ctx context.Context) opentracing.Span {
	if !GetOpentracingEnabled(ctx) {
		return nil
	}
	value := ctx.Value(share.OpentracingSpanServerKey)
	if value != nil {
		span, ok := value.(opentracing.Span)
		if ok {
			return span
		}
		return nil
	}
	return opentracing.SpanFromContext(ctx)
}

// StartOpentracingSpan 从Context里产生一个新的节点
func StartOpentracingSpan(ctx context.Context, name string) (opentracing.Span, context.Context) {
	value := ctx.Value(share.OpentracingSpanServerKey)
	if value != nil {
		span, ok := value.(opentracing.Span)
		if ok {
			spanCtx := opentracing.ContextWithSpan(ctx, span)
			return opentracing.StartSpanFromContext(spanCtx, name)
		}
		return nil, ctx
	}
	if !GetOpentracingEnabled(ctx) {
		return nil, ctx
	}
	//这个是正常的http请求来的
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return opentracing.StartSpanFromContext(ctx, name)
	} else {
		return nil, ctx
	}
}

// TraceHTTPClient 使用链路追踪，如果开启了
func TraceHTTPClient(client *ghttp.Client, req *http.Request) (*ghttp.ClientResponse, error) {
	span, _ := StartOpentracingSpan(req.Context(), "http-client")
	if span != nil {
		span.SetTag("http.url", req.URL.String())
		span.SetTag("http.method", req.Method)
		defer span.Finish()
	}
	return client.Next(req)
}

// Context 注入到http request里面，用于上下文共享
type Context struct {
	User *ContextUser
	Data g.Map
	I18n *I18n
	Span opentracing.Span
}

// ContextUser 在请求上下文中的用户信息
type ContextUser struct {
	UID               uint32
	AppID             AppID //用户对应的APPID
	Salt              string
	Platform          string //用户平台
	Time              uint32 //令牌生成时间
	Agent             string //用户的Agent
	Channel           string //用户所属渠道
	Package           string //当前请求的包名
	Language          string //用户的原始语言
	Area              string //用户的原始地区
	NativeVersion     uint32 //客户端版本号 ip2long
	NativeMainVersion uint32 //build号为0的版本
	JsVersion         uint32 //已经无用
	Mac               string
	DeviceName        string
	Did               string
	IsSimulator       func() bool //是否模拟器
	IsRoot            func() bool //是否Root设备
	IsSigned          func(uint32) bool
}

// GetLanguage 返回当前用户的请求语言
// 优先级 query lang > header USER_LANGUAGE > header Accept-Language
// todo... 可能需要进行归一
func (ctx *ContextUser) GetLanguage() string {
	lang := strings.ToLower(ctx.Language)
	if _, ok := languages[lang]; ok {
		return lang
	}
	return defaultLang
}

// GetChannel 返回当前用户所属渠道，来自用户登录时产生的令牌里
func (ctx *ContextUser) GetChannel() string {
	return ctx.Channel
}

// GetAgent 返回用户所属平台，ios|android|pc
func (ctx *ContextUser) GetAgent() string {
	return ctx.Agent
}

// GetAppID 返回用户是哪个独立APP的
func (ctx *ContextUser) GetAppID() AppID {
	return ctx.AppID
}

// GetAppName 返回当前请求所属APP的名字
// todo...
func (ctx *ContextUser) GetAppName() string {
	return ""
}

// IsLogined 判断当前用户是否登录
func (ctx *ContextUser) IsLogined() bool {
	return ctx.UID > 0 && ctx.AppID > 0 && ctx.Time > 0
}

// IsMinApp 判断是不是小程序，当前请求
func (ctx *ContextUser) IsMinApp() bool {
	return strings.Contains(ctx.Package, ".mini.")
}
