package query

type ReqVoiceLoverMain struct {
}

type ReqAlbumList struct {
	Choice    uint32 `v:"choice@required"` // 0-默认 1-精选 99-专题
	SubjectId uint64 `json:"subject_id"`   // Choice=99的时候，需要传专题id
	Paginator
}

type ReqRecUserList struct {
	Paginator
}

type ReqAlbumDetail struct {
	AlbumId uint64 `v:"album_id@required"` // 专辑id
}

type ReqAlbumComments struct {
	AlbumId uint64 `v:"album_id@required"` // 专辑id
	Paginator
}

type ReqCommentAlbum struct {
	AlbumId uint64 `v:"album_id@required"` // 专辑id
	Comment string `v:"comment@required"`  // 评论内容
}

type ReqAudioDetail struct {
	AudioId uint64 `v:"audio_id@required"` // 音频id
}

type ReqAudioComments struct {
	AudioId uint64 `v:"audio_id@required"` // 音频id
	Type    uint32 `v:"type@required"`     // 0-普通评论 1-弹幕
	Paginator
}

type ReqCommentAudio struct {
	AudioId uint64 `v:"audio_id@required"` // 音频id
	Comment string `v:"comment@required"`  // 评论内容
	Type    uint32 `v:"type@required"`     // 0-普通评论 1-弹幕
}

type ReqCollectVoiceLover struct {
	Id   uint64 `v:"id@required"`   // 资源id
	Type uint32 `v:"type@required"` // 0-专辑 1-音频
	From uint32 `v:"from@required"` // 0-取消收藏 1-收藏
}

type ReqCollectAlbumList struct {
	Paginator
}

type ReqCollectAudioList struct {
	Paginator
}

type ReqVoiceLoverPost struct {
	Resource    string `v:"resource@required"` // 音频资源
	Seconds     uint32 `v:"seconds@required"`  // 音频时长 单位秒
	Title       string `v:"title@required"`    // 标题
	Source      int32  `v:"source@required"`   // 来源 1:原创 2:搬运
	Cover       string `v:"cover@required"`    // 封面
	Desc        string `v:"desc@required"`     // 简介
	EditDub     string `json:"edit_dub"`       // 编辑配音
	EditContent string `json:"edit_content"`   // 编辑文案
	EditPost    string `json:"edit_post"`      // 编辑后期
	EditCover   string `json:"edit_cover"`     // 编辑封面
	Labels      string `json:"labels"`         // 标签
}

type ReqReport struct {
	Id   uint64 `v:"id@required"`   // 资源id
	From uint32 `v:"from@required"` // 0-专辑 1-音频
}

type ReqPlayStatReport struct {
	AlbumId uint64 `v:"album_id@required"` // 专辑id
	AudioId uint64 `v:"audio_id@required"` // 音频id
}

type ReqShareAlbum struct {
	AlbumId uint64 `v:"album_id@required"` // 专辑id
}

type ReqShareAudio struct {
	AudioId uint64 `v:"audio_id@required"` // 音频id
}

type ReqShareAlbumFans struct {
	AlbumId uint64 `v:"album_id@required"` // 专辑id
}

type ReqShareAudioFans struct {
	AudioId uint64 `v:"audio_id@required"` // 音频id
}

type ReqActivityMain struct {
	ActivityId uint32 `v:"activity_id@required"` // 活动id
}
