package voice_lover

type VoiceLoverAudioEsModel struct {
	Id          uint64   `json:"id"`
	PubUid      uint32   `json:"pub_uid"`
	Title       string   `json:"title"`
	Cover       string   `json:"cover"`
	Desc        string   `json:"desc"`
	CreateTime  uint64   `json:"create_time"`
	Labels      []string `json:"labels"`
	Source      int32    `json:"source"`
	AuditStatus int32    `json:"audit_status"`
	Albums      []uint64 `json:"albums"`
	HasAlbum    int32    `json:"has_album"`
	OpUid       uint64   `json:"op_uid"`
	Resource    string   `json:"resource"`
	Seconds     int32    `json:"seconds"`
}
