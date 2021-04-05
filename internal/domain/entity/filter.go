package entity

type FilterOperation string

const (
	Equal FilterOperation = "="
	Like  FilterOperation = "~"
)
