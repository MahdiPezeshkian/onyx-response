// Package response provides transport-independent API response models and
// constructors. JSON field names and omission behavior match the historical
// Onyx.Common/responswrapper contract.
package response

// PaginationRequest describes a client pagination request.
type PaginationRequest struct {
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
	SortOrder  string `json:"sort_order"`
}

// Base holds shared success/error envelope fields.
type Base struct {
	Message      string `json:"message,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	HasError     bool   `json:"has_error"`
}

// Response is a generic API success/error envelope.
type Response[T any] struct {
	Base
	Value *T `json:"value,omitempty"`
}

// Success builds a successful response envelope.
// An empty message is omitted from JSON (omitempty).
// A nil value is omitted from JSON (omitempty).
func Success[T any](value *T, message string) *Response[T] {
	return &Response[T]{
		Base: Base{
			Message:      message,
			ErrorMessage: "",
			HasError:     false,
		},
		Value: value,
	}
}

// Failure builds an error response envelope with a nil value.
// An empty errorMessage is omitted from JSON (omitempty).
// Message is left empty so it is omitted from JSON.
func Failure(errorMessage string) *Response[any] {
	return &Response[any]{
		Base: Base{
			Message:      "",
			ErrorMessage: errorMessage,
			HasError:     true,
		},
		Value: nil,
	}
}

// PagedResponse is a generic paginated API envelope.
// StatusCode is a caller-supplied body field preserved for JSON compatibility;
// this package does not map domain errors to HTTP status codes.
type PagedResponse[T any] struct {
	Base
	Items       []T  `json:"items,omitempty"`
	PageNumber  int  `json:"page_number"`
	PageSize    int  `json:"page_size"`
	ItemsCount  int  `json:"items_count"`
	TotalCount  int  `json:"total_count"`
	HasNextPage bool `json:"has_next_page"`
	StatusCode  int  `json:"status_code"`
}

// PagedSuccess builds a successful paginated response.
//
// Pagination contract: PageNumber is 1-based (first page is 1), matching the
// historical JSON fixtures and common API usage. HasNextPage is:
//
//	(pageNumber-1)*pageSize + len(items) < totalCount
//
// Both nil and empty items slices are omitted from JSON (omitempty).
// StatusCode is a caller-supplied JSON body field for legacy compatibility.
func PagedSuccess[T any](items []T, pageNumber, pageSize, totalCount, statusCode int, message string) *PagedResponse[T] {
	itemsCount := len(items)
	next := computeHasNextPage(pageNumber, pageSize, itemsCount, totalCount)

	return &PagedResponse[T]{
		Base: Base{
			Message:      message,
			ErrorMessage: "",
			HasError:     false,
		},
		Items:       items,
		PageNumber:  pageNumber,
		PageSize:    pageSize,
		ItemsCount:  itemsCount,
		TotalCount:  totalCount,
		HasNextPage: next,
		StatusCode:  statusCode,
	}
}

// PagedFailure builds a paginated error response with zeroed paging fields.
func PagedFailure[T any](statusCode int, errorMessage string) *PagedResponse[T] {
	return &PagedResponse[T]{
		Base: Base{
			Message:      "",
			ErrorMessage: errorMessage,
			HasError:     true,
		},
		Items:       nil,
		PageNumber:  0,
		PageSize:    0,
		ItemsCount:  0,
		TotalCount:  0,
		HasNextPage: false,
		StatusCode:  statusCode,
	}
}

// computeHasNextPage implements the 1-based pagination contract.
// Non-positive pageNumber or pageSize yields false.
func computeHasNextPage(pageNumber, pageSize, itemsCount, totalCount int) bool {
	if pageNumber < 1 || pageSize < 1 {
		return false
	}
	consumed := (pageNumber-1)*pageSize + itemsCount
	return consumed < totalCount
}
