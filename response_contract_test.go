package response_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/MahdiPezeshkian/onyx-response"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func mustFixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", "fixtures", name)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return b
}

func assertJSONEqual(t *testing.T, want []byte, got any) {
	t.Helper()
	gotBytes, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	gotBytes = append(gotBytes, '\n')

	var wantNorm, gotNorm any
	if err := json.Unmarshal(want, &wantNorm); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	if err := json.Unmarshal(gotBytes, &gotNorm); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}
	wantCanon, _ := json.MarshalIndent(wantNorm, "", "  ")
	gotCanon, _ := json.MarshalIndent(gotNorm, "", "  ")
	if !bytes.Equal(wantCanon, gotCanon) {
		t.Fatalf("JSON mismatch\nwant:\n%s\ngot:\n%s", wantCanon, gotCanon)
	}
}

func TestSuccessContract_WithValue(t *testing.T) {
	user := &User{ID: 1, Name: "Ada"}
	assertJSONEqual(t, mustFixture(t, "success_with_value.json"), response.Success(user, "token refreshed successfully"))
}

func TestSuccessContract_NilValue(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "success_nil_value.json"), response.Success[User](nil, "ok"))
}

func TestSuccessContract_EmptyMessage(t *testing.T) {
	user := &User{ID: 1, Name: "Ada"}
	assertJSONEqual(t, mustFixture(t, "success_empty_message.json"), response.Success(user, ""))
}

func TestFailureContract_Basic(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "error_basic.json"), response.Failure("invalid refresh token"))
}

func TestFailureContract_EmptyMessage(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "error_empty_message.json"), response.Failure(""))
}

func TestPagedSuccessContract(t *testing.T) {
	users := []User{{ID: 1, Name: "Ada"}}
	assertJSONEqual(t, mustFixture(t, "paged_success.json"), response.PagedSuccess(users, 1, 10, 25, 200, "ok"))
}

func TestPagedSuccessContract_ExactPageNoNext(t *testing.T) {
	users := []User{{ID: 1, Name: "Ada"}}
	assertJSONEqual(t, mustFixture(t, "paged_success_exact_page.json"), response.PagedSuccess(users, 3, 10, 21, 200, "ok"))
}

func TestPagedSuccessContract_NilItems(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "paged_nil_items.json"), response.PagedSuccess[User](nil, 1, 10, 0, 200, "ok"))
}

func TestPagedSuccessContract_EmptyItems(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "paged_empty_items.json"), response.PagedSuccess([]User{}, 1, 10, 0, 200, "ok"))
}

func TestPagedFailureContract(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "paged_error.json"), response.PagedFailure[User](404, "not found"))
}

func TestPagedFailureContract_EmptyMessage(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "paged_error_empty.json"), response.PagedFailure[User](500, ""))
}

func TestPaginationRequestContract(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "pagination_request.json"), response.PaginationRequest{
		PageNumber: 2,
		PageSize:   20,
		SortOrder:  "asc",
	})
}

func TestPaginationRequestContract_Empty(t *testing.T) {
	assertJSONEqual(t, mustFixture(t, "pagination_request_empty.json"), response.PaginationRequest{})
}

// TestLegacyQuirk_NilAndEmptyItemsBothOmitField documents that both nil and
// empty slices omit "items" under omitempty — a preserved legacy behavior.
func TestLegacyQuirk_NilAndEmptyItemsBothOmitField(t *testing.T) {
	nilJSON, err := json.Marshal(response.PagedSuccess[User](nil, 1, 10, 0, 200, "ok"))
	if err != nil {
		t.Fatal(err)
	}
	emptyJSON, err := json.Marshal(response.PagedSuccess([]User{}, 1, 10, 0, 200, "ok"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(nilJSON, emptyJSON) {
		t.Fatalf("nil and empty items must serialize identically; nil=%s empty=%s", nilJSON, emptyJSON)
	}
	if bytes.Contains(nilJSON, []byte(`"items"`)) {
		t.Fatalf("items must be omitted for nil/empty slices: %s", nilJSON)
	}
}
