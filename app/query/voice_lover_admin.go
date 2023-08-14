package query

type ReqAdminVoiceLoverAudioList struct {
	StartTime   uint64
	EndTime     uint64
	UserStr     string
	Source      int32
	Label       string
	Order       string
	AuditStatus int32 // -1:全部 0:未审 1:通过 2:未通过
	Paginator
}

type ReqAdminVoiceLoverAudioDetail struct {
	Id uint64
}

type ReqAdminVoiceLoverAudioUpdate struct {
	Id     uint64
	Title  string
	Desc   string
	Labels string
	OpUid  uint64
}

type ReqAdminVoiceLoverAudioAudit struct {
	Id          uint64
	AuditStatus int32 // 1:通过 2:不通过
	AuditReason int32 //
	OpUid       uint64
}

type ReqAdminVoiceLoverAlbumCreate struct {
	Name  string
	Intro string
	Cover string
	OpUid uint64
}

type ReqAdminVoiceLoverAlbumDel struct {
	Id    uint64
	OpUid uint64
}

type ReqAdminVoiceLoverAlbumUpdate struct {
	Id    uint64
	Name  string
	Intro string
	Cover string
	OpUid uint64
}

type ReqAdminVoiceLoverAlbumList struct {
	Paginator
	Name      string
	StartTime uint64
	EndTime   uint64
}

type ReqAdminVoiceLoverAlbumDetail struct {
	Id uint64
}

type ReqAdminVoiceLoverAudioCollectList struct {
	UserStr       string
	AlbumStr      string
	AudioStr      string
	Label         string
	CollectStatus int32 //-1:全部 0:未收录 1:已收录
	Order         string
	Paginator
}

type ReqAdminVoiceLoverAudioCollect struct {
	AudioId uint64
	AlbumId uint64
	Type    int32 //0:收录 1:移除
}
