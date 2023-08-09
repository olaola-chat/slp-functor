package query

type Paginator struct {
	Page  uint32 `v:"page@integer"`                // 第几页, 范围[1,10]
	Limit uint32 `v:"limit@integer|min:10|max:50"` // 每页多少, 范围[10,50]
}
