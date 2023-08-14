package dao

type voiceLoverUserCollectDao struct {
}

var VoiceLoverUserCollectDao = &voiceLoverUserCollectDao{}

const (
	CollectTypeAlbum = 0 // 专辑收藏
	CollectTypeAudio = 1 // 音频收藏
)
