package query

type ReqVoiceLoverMain struct {
}

type ReqAlbumList struct {
	Choice    uint32 `v:"choice@required"` // 0-默认 1-精选 2-专题
	SubjectId uint64 `v:"subject_id"`      // Choice=2的时候，需要传专题id
	Paginator
}

type ReqRecUserList struct {
	Paginator
}

type ReqAlbumDetail struct {
	AlbumId uint32 `v:"album_id@required"` // 专辑id
}

type ReqVoiceLoverPost struct {
	Resource    string //音频资源
	Seconds     int32  //音频时长 单位秒
	Title       string //标题
	Source      int32  //来源 1:原创 2:搬运
	Cover       string //封面
	Desc        string //简介
	EditDub     string // 编辑配音
	EditContent string //编辑文案
	EditPost    string //编辑后期
	EditCover   string //编辑封面
	Labels      string //标签
}
