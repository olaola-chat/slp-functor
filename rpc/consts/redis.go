package consts

import (
	"fmt"
	"time"
)

const (
	RedisDefault = "default"
)

type RedisKey struct {
	key string
	ttl time.Duration
}

func (k *RedisKey) Key(values ...interface{}) string {
	return fmt.Sprintf(k.key, values...)
}

func (k *RedisKey) Ttl() time.Duration {
	return k.ttl
}

/*
 * 以下为rpc模块的缓存定义 分模块
 */
// 用户收藏模块
var (
	UserCollectAlbumKey = &RedisKey{key: "user.collect.album.%d_%d", ttl: 7 * 24 * time.Hour} // %d_%d = uid_albumId
	UserCollectVoiceKey = &RedisKey{key: "user.collect.voice.%d_%d", ttl: 7 * 24 * time.Hour} // %d_%d = uid_voiceId
)
