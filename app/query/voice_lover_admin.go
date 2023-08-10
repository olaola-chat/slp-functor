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
