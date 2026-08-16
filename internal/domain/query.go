package domain

type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

type QueryOptions struct {
	Page    int
	Limit   int
	Search  string
	Sort    string
	Filters []Filter
}
