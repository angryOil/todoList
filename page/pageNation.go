package page

var minPage = 0
var minSize = 10
var maxSize = 50

type ReqPage struct {
	Page int
	Size int
}

func NewReqPage(page int, size int) ReqPage {
	rp := ReqPage{}
	if page <= minPage {
		rp.Page = 0
	} else {
		rp.Page = page
	}
	if size <= minSize {
		rp.Size = minSize
	} else if size >= maxSize {
		rp.Size = maxSize
	} else {
		rp.Size = size
	}
	return rp
}

type Pagination[T any] struct {
	Contents    []T `json:"contents"`
	Total       int `json:"total_content"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
}

// page 는 0번째 page 부터 시작합니다.

func GetPagination[T any](contents []T, rp ReqPage, totalCount int) Pagination[T] {

	return Pagination[T]{
		Contents:    contents,
		Total:       totalCount,
		CurrentPage: rp.Page,
		LastPage:    totalCount / rp.Size,
	}
}
