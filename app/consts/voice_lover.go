package consts

var AuditAudioReasonMap = map[uint64]string{
	1: "音频内容不合规",
	2: "封面不合规",
	3: "音频标题不合规",
	4: "音频简介不合规",
	5: "疑似非原创",
	6: "内容质量不符合要求",
}

const (
	AudioPassTopic = "audio_pass"
)
