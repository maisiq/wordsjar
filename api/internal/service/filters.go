package service

type field string

const (
	FieldID field = "id"
	FieldEN field = "en"
)

type QueryParams struct {
	Limit      uint8
	SortBy     field
	Desc       bool
	Pagination *Pagination
}

type UserWordsFilter struct {
	TestMode bool
	WordID   string
}

type Filter func(*UserWordsFilter)

func WithTestMode() Filter {
	return func(o *UserWordsFilter) {
		o.TestMode = true
	}
}

func WithWordID(id string) Filter {
	return func(o *UserWordsFilter) {
		o.WordID = id
	}
}
