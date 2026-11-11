package response_test

import (
	"testing"

	"github.com/MahdiPezeshkian/onyx-response"
)

func TestPagedHasNextBoundaries(t *testing.T) {
	type row struct {
		page, size, total int
		items             int
		wantNext          bool
		name              string
	}
	cases := []row{
		{1, 10, 25, 10, true, "first full page with more"},
		{1, 10, 10, 10, false, "exact single page"},
		{1, 10, 0, 0, false, "empty total"},
		{2, 10, 25, 10, true, "middle page"},
		{3, 10, 25, 5, false, "last partial page"},
		{3, 10, 21, 1, false, "legacy exact fixture shape"},
		{1, 10, 20, 10, true, "first page of two"},
		{2, 10, 20, 10, false, "second exact page"},
		{0, 10, 100, 10, false, "invalid page 0"},
		{1, 0, 100, 0, false, "invalid page size"},
		{-1, 10, 100, 10, false, "negative page"},
	}
	for _, tc := range cases {
		items := make([]int, tc.items)
		got := response.PagedSuccess(items, tc.page, tc.size, tc.total, 200, "ok")
		if got.HasNextPage != tc.wantNext {
			t.Fatalf("%s: HasNextPage=%v want %v (page=%d size=%d items=%d total=%d)",
				tc.name, got.HasNextPage, tc.wantNext, tc.page, tc.size, tc.items, tc.total)
		}
		if got.ItemsCount != tc.items || got.PageNumber != tc.page || got.PageSize != tc.size {
			t.Fatalf("%s: metadata mismatch %#v", tc.name, got)
		}
	}
}
