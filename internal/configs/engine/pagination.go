package engine

import "math"

type Pagination struct {
	Pages    int `json:"pages"`
	Total    int `json:"total"`
	PageSize int `json:"pageSize"`
}

func NewPagination(total, size, currentLimit int) Pagination {
	pag := Pagination{
		Total:    total,
		PageSize: size,
		Pages:    int(math.Ceil(float64(total) / float64(currentLimit))),
	}
	return pag
}
