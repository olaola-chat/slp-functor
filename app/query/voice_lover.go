package query

type ReqVoiceLoverMain struct {
	Uid uint64
}

type ReqVoiceLoverPost struct {
	Resource    string //音频资源
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
