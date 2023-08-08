package query

type ReqVoiceLoverMain struct {
}

type ReqAlbumList struct {
	Paginator
}

type ReqVoiceLoverPost struct {
	Resource    string `json:"resource"`     //音频资源
	Title       string `json:"tile"`         //标题
	Source      int32  `json:"source"`       //来源 1:原创 2:搬运
	Cover       string `json:"cover"`        //封面
	Desc        string `json:"desc"`         //简介
	EditDub     string `json:"edit_dub"`     // 编辑配音
	EditContent string `json:"edit_content"` //编辑文案
	EditPost    string `json:"edit_post"`    //编辑后期
	EditCover   string `json:"edit_cover"`   //编辑封面
	Labels      string `json:"labels"`       //标签
}
