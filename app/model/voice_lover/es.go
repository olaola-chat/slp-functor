package voice_lover

type VoiceLoverEsModel struct {
	Id          uint64   `json:"id"`
	PubUid      uint64   `json:"pub_uid"`
	Title       string   `json:"title"`
	Cover       string   `json:"cover"`
	Desc        string   `json:"desc"`
	CreateTime  uint64   `json:"create_time"`
	Labels      []string `json:"labels"`
	Source      uint32   `json:"source"`
	AuditStatus uint32   `json:"audit_status"`
	Albums      []uint64 `json:"albums"`
	HasAlbum    int32    `json:"has_album"`
	OpUid       uint64   `json:"op_uid"`
}
