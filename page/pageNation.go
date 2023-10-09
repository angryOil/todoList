package page

var minPage = 0
var minSize = 10
var maxSize = 50

type ReqPage struct {
	page int
	size int
}

func NewReqPage(page int, size int) ReqPage {
	rp := ReqPage{}
	if page <= minPage {
		rp.page = 0
	} else {
		rp.page = page
	}
	if size <= minSize {
		rp.size = minSize
	} else if size >= maxSize {
		rp.size = maxSize
	} else {
		rp.size = size
	}
	return rp
}

type Pagination[T any] struct {
	Contents    []T `json:"contents"`
	Total       int `json:"total"`
	CurrentPage int `json:"current"`
	LastPage    int `json:"last"`
}

func GetPagination[T any](contents []T, rp ReqPage, currentPage int, totalCount int) Pagination[T] {
	return Pagination[T]{
		Contents:    contents,
		Total:       totalCount,
		CurrentPage: currentPage + 1,
		LastPage:    (totalCount / rp.size) + 1,
	}
}
