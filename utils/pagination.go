package utils

import "strconv"

type Pagination struct {
	Page  string `json:"page"`
	Limit string `json:"limit"`
	Size  string `json:"size"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetPage() int {
	page, _ := strconv.Atoi(p.Page)
	if page == 0 {
		return 1
	}
	return page
}

func (p *Pagination) GetLimit() int {
	limit, _ := strconv.Atoi(p.Limit)
	if limit == 0 {
		return 10
	}
	return limit
}

func (p *Pagination) GetSize() int {
	size, _ := strconv.Atoi(p.Size)
	return size
}

func (p *Pagination) PaginateQuery() string {
	return " LIMIT " + strconv.Itoa(p.GetLimit()) + " OFFSET " + strconv.Itoa(p.GetOffset())
}
