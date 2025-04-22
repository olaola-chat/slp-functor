package nsq

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"sync/atomic"
	"time"

	"github.com/olaola-chat/slp-functor/library/tool"

	"github.com/gogf/gf/frame/gins"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

var topics map[string]string = map[string]string{
	"csms.nsq": NsqGroupCircle,
}

const (
	// NsqGroupDefault 在config中对应的配置名字
	NsqGroupDefault = "default"
	NsqGroupCircle  = "circle"
	// NsqGroup 在config中对应的配置名字
	NsqGroup = "nsq"
	// NsqConfigNsdName 在config中对应的配置名字
	NsqConfigNsdName = "go-nsq"
)

// GetNsdGroup 根据topic来返回对应的生产者对象
func GetNsdGroup(topic string) ([]*nsq.Producer, error) {
	groupName := ""
	ok := false
	// 只读map，没问题，直接访问
	if groupName, ok = topics[topic]; !ok {
		// 直接崩掉，不能掩盖问题
		panic(gerror.New("nsq topic not defined" + topic))
	}

	instanceKey := fmt.Sprintf("self-go-nsq-group.%s", groupName)
	result := gins.GetOrSetFuncLock(instanceKey, func() interface{} {
		// 因配置错误导致的，需要直接panic
		addrs := g.Cfg().GetStrings(fmt.Sprintf("%s.%s", NsqConfigNsdName, groupName))
		if len(addrs) == 0 {
			panic(gerror.New("nsq config name error"))
		}

		ps := []*nsq.Producer{}
		ip, err := tool.IP.LocalIPv4s()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			cfg := nsq.NewConfig()
			cfg.ReadTimeout = time.Second * 30
			cfg.HeartbeatInterval = time.Second * 15
			cfg.Hostname = ip
			producer, err := nsq.NewProducer(addr, cfg)
			if err != nil {
				panic(gerror.Wrap(err, "nsq create producer error"))
			}
			ps = append(ps, producer)
		}

		return ps
	})

	if result == nil {
		return nil, gerror.New("error get nsd")
	}

	if nsds, ok := result.([]*nsq.Producer); ok {
		return nsds, nil
	}

	return nil, gerror.New("error get nsd")
}

// NewNsqClient 实例化一个Client对象
func NewNsqClient() *Client {
	return &Client{}
}

// Client 定义消息发送对象
type Client struct {
	index uint64
}

// Send 发送消息，使用json编码
// todo... 兼容php现在的代码
func (c *Client) Send(topic string, body interface{}, delay ...time.Duration) error {
	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	nsds, err := GetNsdGroup(topic)
	if err != nil {
		return err
	}
	index := atomic.AddUint64(&c.index, 1)
	client := nsds[int(index%uint64(len(nsds)))]
	if len(delay) > 0 {
		err = client.DeferredPublish(topic, delay[0], bytes)
	} else {
		err = client.Publish(topic, bytes)
	}
	return err
}

func (c *Client) SendBytes(topic string, body []byte, delay ...time.Duration) error {
	nsds, err := GetNsdGroup(topic)
	if err != nil {
		return err
	}
	index := atomic.AddUint64(&c.index, 1)
	client := nsds[int(index%uint64(len(nsds)))]
	if len(delay) > 0 {
		err = client.DeferredPublish(topic, delay[0], body)
	} else {
		err = client.Publish(topic, body)
	}
	return err
}

// SendIgoreError 当前不想被ci警告时，使用这个
// 需要考虑清楚
func (c *Client) SendIgoreError(topic string, body interface{}, delay ...time.Duration) {
	err := c.Send(topic, body, delay...)
	if err != nil {
		g.Log().Println(err)
	}
}

// 多发
func (c *Client) MultiSend(topic string, body [][]byte) error {
	nsds, err := GetNsdGroup(topic)
	if err != nil {
		return err
	}
	index := atomic.AddUint64(&c.index, 1)
	client := nsds[int(index%uint64(len(nsds)))]
	err = client.MultiPublish(topic, body)
	return err
}
