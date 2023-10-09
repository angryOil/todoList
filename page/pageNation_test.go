package page

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPagination(t *testing.T) {
	rp := NewReqPage(0, 0)

	p := GetPagination(make([]int, 0), rp, 0, 90)
	assert.Equal(t, 10, p.LastPage)

	p = GetPagination(make([]int, 0), rp, 0, 67)
	assert.Equal(t, 7, p.LastPage)

	p = GetPagination(make([]int, 0), rp, 0, 1)
	assert.Equal(t, 1, p.LastPage)

	p = GetPagination(make([]int, 0), rp, 0, 0)
	assert.Equal(t, 1, p.LastPage)

	p = GetPagination(make([]int, 0), rp, 11, 11)
	assert.Equal(t, 2, p.LastPage)

}
