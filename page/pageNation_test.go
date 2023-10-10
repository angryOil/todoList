package page

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPagination(t *testing.T) {
	rp := NewReqPage(0, 0)

	p := GetPagination(make([]int, 0), rp, 90)
	assert.Equal(t, 9, p.LastPage)

	p = GetPagination(make([]int, 0), rp, 67)
	assert.Equal(t, 6, p.LastPage)

	p = GetPagination(make([]int, 0), rp, 1)
	assert.Equal(t, 0, p.LastPage)

	p = GetPagination(make([]int, 0), rp, 0)
	assert.Equal(t, 0, p.LastPage)

	p = GetPagination(make([]int, 0), rp, 11)
	assert.Equal(t, 1, p.LastPage)

}
