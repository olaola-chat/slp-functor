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

// 声音恋人音频模块
var (
	// 播放次数 %d = id
	VoiceLoverAudioPlayCount = &RedisKey{key: "vl.audio.play_count.%d", ttl: -1}
)

// 用户收藏模块
var (
	// 声音恋人收藏专辑 %d_%d = uid_albumId
	UserCollectAlbumKey = &RedisKey{key: "user.collect.album.%d_%d", ttl: 7 * 24 * time.Hour}
	// 声音恋人收藏音频 %d_%d = uid_voiceId
	UserCollectAudioKey = &RedisKey{key: "user.collect.audio.%d_%d", ttl: 7 * 24 * time.Hour}
)
