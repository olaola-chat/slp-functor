package library

import (
	"time"

	cnsq "github.com/olaola-chat/rbp-functor/library/nsq"
	"github.com/olaola-chat/rbp-functor/library/tool"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	nsq "github.com/nsqio/go-nsq"
)

const (
	//NsqGroupDefault 默认的NSD机器组配置名字
	NsqGroupDefault = cnsq.NsqGroupDefault
	//NsqConfigLookupName 集群的lookup配置名字
	NsqConfigLookupName = "go-nsq.lookup"
	NsqGroupCircle      = cnsq.NsqGroupCircle
)

// NsqClient 返回cnsq client 对象
func NsqClient() *cnsq.Client {
	return cnsq.NewNsqClient()
}

// NewNsqWorker 实例化 NsqWorker
func NewNsqWorker(topic, channel string, handler NsqHandleMessage) *NsqWorker {
	return &NsqWorker{
		topic:   topic,
		channel: channel,
		handler: handler,
	}
}

// NsqHandleMessage 定义NewNsqWorker的回调方式
type NsqHandleMessage func(message *nsq.Message) error

// NsqWorker 基于topic消费nsd
type NsqWorker struct {
	client  *nsq.Consumer
	topic   string
	channel string
	handler NsqHandleMessage
}

// HandleMessage nsq 基类的回调接口
func (s *NsqWorker) HandleMessage(message *nsq.Message) error {
	//反序列化数据为JSON
	return s.handler(message)
}

// Connect 建立连接函数
func (s *NsqWorker) Connect() error {
	var err error
	ip, _ := tool.IP.LocalIPv4s()
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second * 30 //设置重连时间
	cfg.HeartbeatInterval = time.Second * 5
	cfg.Hostname = ip
	nsds, err := cnsq.GetNsdGroup(s.topic)
	if err != nil {
		return err
	}
	cfg.MaxInFlight = len(nsds)
	g.Log().Infof("nsq.Consumer MaxInFlight=%d", cfg.MaxInFlight)
	s.client, err = nsq.NewConsumer(s.topic, s.channel, cfg) // 新建一个消费者
	if err != nil {
		return err
	}
	s.client.SetLogger(nil, 0) //屏蔽系统日志
	s.client.AddHandler(s)     // 添加消费者接口
	s.client.SetBehaviorDelegate(s)
	//建立NSQLookupd连接
	config := g.Cfg().GetStrings(NsqConfigLookupName)
	g.Log().Infof("nsq.NsqConfigLookupName %v\n", config)
	if len(config) == 0 {
		panic(gerror.New("empty lookup for nsq"))
	}
	if err := s.client.ConnectToNSQLookupds(config); err != nil {
		g.Log().Infof("nsq.NsqConfigLookupName  err = %s\n", err.Error())
		return err
	}
	g.Log().Infof("nsq.NsqConfigLookupName  ok\n")

	return nil
}

// Stop 停止
func (s *NsqWorker) Stop() {
	s.client.Stop()
}

// Stats 统计数据
func (s *NsqWorker) Stats() *nsq.ConsumerStats {
	return s.client.Stats()
}

// DiscoveryFilter 对现在线上的NSD ip 进行转换
// todo... 伴伴环境需要填写具体的转化IP
func (s NsqWorker) Filter(addrs []string) []string {
	values := []string{}
	replace := map[string]string{
		"10.80.153.177:4150": "172.16.0.179:4250",
		"10.31.52.45:4150":   "172.16.0.179:4050",
		"10.31.52.45:4152":   "172.16.0.179:4052",
		"10.31.52.45:4154":   "172.16.0.179:4054",
		"10.31.52.45:4156":   "172.16.0.179:4056",
		"10.31.52.45:4158":   "172.16.0.179:4058",
		"10.81.45.178:4150":  "172.16.0.179:4150",
		"10.81.45.178:4152":  "172.16.0.179:4152",
		"10.81.45.178:4154":  "172.16.0.179:4154",
		"10.81.45.178:4156":  "172.16.0.179:4156",
	}
	for i := 0; i < len(addrs); i++ {
		addr := addrs[i]
		if val, ok := replace[addr]; ok {
			values = append(values, val)
		} else {
			values = append(values, addr)
		}
	}
	return values
}

// Connect 建立连接函数(可以指定并发消费的消息数)
func (s *NsqWorker) ConnectWithConcurrency(cNum int) error {
	var err error
	ip, _ := tool.IP.LocalIPv4s()
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second * 30 //设置重连时间
	cfg.HeartbeatInterval = time.Second * 5
	cfg.Hostname = ip
	nsds, err := cnsq.GetNsdGroup(s.topic)
	if err != nil {
		return err
	}
	cfg.MaxInFlight = len(nsds)
	if cfg.MaxInFlight < cNum {
		cfg.MaxInFlight = cNum
	}
	g.Log().Infof("nsq.Consumer MaxInFlight=%d", cfg.MaxInFlight)
	s.client, err = nsq.NewConsumer(s.topic, s.channel, cfg) // 新建一个消费者
	if err != nil {
		return err
	}
	s.client.SetLogger(nil, 0)              //屏蔽系统日志
	s.client.AddConcurrentHandlers(s, cNum) // 添加消费者接口
	s.client.SetBehaviorDelegate(s)
	//建立NSQLookupd连接
	config := g.Cfg().GetStrings(NsqConfigLookupName)
	g.Log().Infof("nsq.NsqConfigLookupName %v\n", config)
	if len(config) == 0 {
		panic(gerror.New("empty lookup for nsq"))
	}
	if err := s.client.ConnectToNSQLookupds(config); err != nil {
		g.Log().Infof("nsq.NsqConfigLookupName  err = %s\n", err.Error())
		return err
	}
	g.Log().Infof("nsq.NsqConfigLookupName  ok\n")

	return nil
}
