# onyx-response

[![Go Version](https://img.shields.io/badge/go-1.27-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Transport-independent API response models and constructors. Extracted from
`Onyx.Common/responswrapper` with a preserved JSON contract.

## Installation

```bash
go get github.com/MahdiPezeshkian/onyx-response@v0.1.0
```

## Import path

```go
import "github.com/MahdiPezeshkian/onyx-response"
```

Package name: `response`

## Basic usage

```go
ok := response.Success(&user, "token refreshed successfully")
fail := response.Failure("invalid refresh token")
page := response.PagedSuccess(users, 1, 10, 25, 200, "ok")
```

Example success JSON:

```json
{
  "message": "token refreshed successfully",
  "has_error": false,
  "value": { "id": 1, "name": "Ada" }
}
```

Example error JSON:

```json
{
  "error_message": "invalid refresh token",
  "has_error": true
}
```

## Public API overview

| API | Purpose |
|---|---|
| `Success` / `Failure` | Generic success/error envelopes (`value`, `has_error`, …) |
| `PagedSuccess` / `PagedFailure` | Paginated envelopes (`items`, paging fields, `status_code`) |
| `PaginationRequest` | Client pagination request model |
| `Base`, `Response`, `PagedResponse` | Exported structs with stable JSON tags |

Constructor names were simplified from the historical `NewRestApiResponseV2`
style; JSON output is covered by golden fixtures under `testdata/fixtures`.

## Design decisions

- Transport-independent: no `net/http`, Gin, Echo, or Fiber.
- Does not map domain errors to HTTP status codes.
- `PagedResponse.StatusCode` is a **caller-supplied JSON body field** preserved
  for compatibility; it is not interpreted by this package.
- Pagination is **1-based** (`page_number` starts at 1). `has_next_page` uses
  `(pageNumber-1)*pageSize + itemsCount < totalCount`.
- Empty `message` / `error_message` and nil `value` are omitted via `omitempty`.
- Both nil and empty `items` slices omit the `items` field (legacy quirk).

## Dependencies

None beyond the Go standard library.

## Limitations

- JSON contract compatibility is verified with golden fixtures derived from
  `Onyx.Common/responswrapper`. Real frontend integration was not exercised in
  this extraction task.
- No HTTP writing helpers (see `onyx-http`).

## License

MIT — see [LICENSE](LICENSE).
