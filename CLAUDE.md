# blueticks-go regeneration conventions

This file is the authoritative spec for a subagent regenerating the Go SDK
from `openapi.json`. The subagent sees this file, the OpenAPI spec, and the
current contents of the regenerated paths — nothing else. Everything it
produces must follow these rules.

## 1. Boundaries

### You MAY write:

- `<resource>.go` — one file per top-level resource (e.g., `account.go`, `messages.go`)
- `<resource>_test.go` — one file per resource
- `types.go` — **only** the type definitions for response/request shapes (the
  pagination-page structs and the resource DTO structs); see Section 5 for the
  carve-out structure within this file.

### You MAY write (only, with narrow constraints):

- `client.go` — **only** to keep the `Client` struct's resource fields and the
  `NewClient` initializer's resource attachments in sync. Specifically:
  - Add or remove `<Resource> *<Resource>Resource` fields on `Client`.
  - Add or remove the corresponding `c.<Resource> = &<Resource>Resource{client: c}` line in `NewClient`.
  - Do NOT modify `apiKey`/`baseURL`/`httpClient` plumbing, `do()`, retry logic, or `Option` functions.

### You MUST NOT touch:

- `errors.go`
- `transport.go`
- `version.go`
- `webhooks.go`
- `webhooks_verify.go`
- `testhelpers_test.go`
- `webhooks_test.go`
- `webhooks_verify_test.go`
- `go.mod`, `go.sum`
- `examples/**`
- `CHANGELOG.md`, `README.md`, `LICENSE`
- `.github/workflows/**`
- `CLAUDE.md` (this file)

## Spec coverage audit (REQUIRED before completion)

Before producing your final diff, you MUST perform this audit. A regenerated SDK that's missing an op or resource that's in the spec is a P0 bug — it's invisible to every downstream caller until someone notices in production. Do not skip this step.

### Step 1: enumerate every spec operation

Read every path in `openapi.json` under `paths` and every method (`get`/`post`/`put`/`patch`/`delete`) for each path. List them in the form `<HTTP_VERB> <path>` (e.g., `POST /v1/newsletters`).

### Step 2: map each operation to the expected SDK method

Using the Method-name mapping table elsewhere in this document, derive the expected `<resource>.<method>` for each spec op. The resource is the first path segment after `/v1/`. The method name follows the verb table.

### Step 3: confirm coverage

For each `<resource>.<method>` you derived:

- If the resource file does NOT exist (e.g. spec has `/v1/newsletters` but no resource file), **CREATE the file** following the structure below + wire it into the SDK entrypoint (create `<resource>.go` at module root, add the resource field to `Client` struct in `client.go` (only inside REGEN-BOUNDARY regions)).
- If the resource file exists but lacks the method, **ADD the method**.
- If a method exists for an operation that's NOT in the spec, **delete it** (the backend has removed the endpoint).

### Step 4: report

In your final summary, include a line:

```
Coverage audit: enumerated N spec ops; M already covered; K added; L deleted; 0 missing.
```

The trailing `0 missing` is the gate. If the audit ends with anything other than `0 missing`, fix it before you finish.

## 2. Types (OpenAPI schema → Go struct)

All response / request DTOs live in **one file**, `types.go`. Group them by
resource with a banner comment between each section:

```go
// --- Account -------------------------------------------------------------

// Account is the response from /v1/account.
type Account struct {
    ID        string  `json:"id"`
    Name      string  `json:"name"`
    Timezone  *string `json:"timezone,omitempty"`
    CreatedAt string  `json:"created_at"`
}
```

### Type mapping

| OpenAPI                             | Go                          |
| ----------------------------------- | --------------------------- |
| `string`                            | `string`                    |
| `string, format: date-time`         | `string`  *(do not parse)*  |
| `integer`                           | `int`                       |
| `number`                            | `float64`                   |
| `boolean`                           | `bool`                      |
| `array, items: <T>`                 | `[]<T>`                     |
| `object` (inline)                   | nested struct in `types.go` |
| `enum: [a, b, c]`                   | `string` (no Go enum type — caller compares to documented constants in a `const ()` block at the top of the section) |
| `nullable: true`                    | pointer (`*T`), tag `,omitempty` |

### Field naming

- Struct field: Go convention (PascalCase, idiomatic acronyms — `URL`, `ID`, not `Url`, `Id`).
- JSON tag: spec wire name verbatim (snake_case).
- Nullable scalars and optional request inputs use `*T` so the JSON encoder distinguishes "absent" from "zero value". For non-nullable required fields, use the value type directly without `omitempty`.

### Pagination wrappers

One per resource. Naming: `<Resource>Page`.

```go
// MessagesPage is the paginated response from MessagesResource.List.
type MessagesPage struct {
    Data       []Message `json:"data"`
    HasMore    bool      `json:"has_more"`
    NextCursor *string   `json:"next_cursor,omitempty"`
}
```

Go pre-1.18 (which this repo targets for compatibility) has no generic page
envelope — each resource gets its own typed page struct.

## 3. Resources (OpenAPI operation → Go method)

One file per path group, named `<resource>.go`. Group by the first path
segment after `/v1/`.

Each file defines one struct + methods:

```go
package blueticks

import (
    "context"
    "net/http"
)

// AccountResource exposes /v1/account.
type AccountResource struct {
    client *Client
}

// Retrieve returns the account associated with the API key used for the request.
func (r *AccountResource) Retrieve(ctx context.Context) (*Account, error) {
    body, err := r.client.do(ctx, requestOptions{method: http.MethodGet, path: "/v1/account"})
    if err != nil {
        return nil, err
    }
    out := &Account{}
    if err := decodeJSON(body, out); err != nil {
        return nil, err
    }
    return out, nil
}
```

### Method-name mapping (from the Feathers operation method)

| Feathers | Go               | When                                    |
| -------- | ---------------- | --------------------------------------- |
| `find`   | `List`           | response is paginated (top-level `data` array + `pagination`) |
| `find`   | `Retrieve`       | response is a single object             |
| `get`    | `Retrieve`       | always                                  |
| `create` | `Create`         | always                                  |
| `patch`  | `Update`         | always                                  |
| `remove` | `Delete`         | always                                  |

### Method signatures

- `ctx context.Context` is **always** the first parameter.
- Path parameters: positional after `ctx`, in path order.
- Query params + body: a single `<Operation>Params` struct (defined in `types.go`)
  passed by value. Optional fields use `*T` so the SDK can distinguish "leave
  default" from "zero".
- Return type: `(*<TypedResponse>, error)` for typed responses, `(*<Resource>Page, error)` for list endpoints, `error` only when the endpoint returns no body (rare — even DELETE routes through `do()`).
- Body: every method routes through `r.client.do(ctx, requestOptions{...})`. Never call `http.Client` directly.
- Doc comment above each method: the operation's `summary` (first line) followed by `\n\n` + `description`. If `summary` is missing, use the first sentence of `description`. If both are missing, write `// <Method> returns the <resource> ...` boilerplate.

### Query / body assembly

`requestOptions` accepts `query url.Values` and `body any` (encoded to JSON).
For `*<Operation>Params`-style inputs, build query.Values in-method and pass
the params struct as `body` for POST/PATCH:

```go
opts := requestOptions{method: http.MethodGet, path: "/v1/messages"}
if params.Cursor != nil {
    opts.query = url.Values{"cursor": []string{*params.Cursor}}
}
```

## 4. Tests

One file per resource, named `<resource>_test.go`.

### Per method: two tests

**Happy path.** Use `newTestServer(t, handler)` from `testhelpers_test.go`.
Inside the handler assert `r.Method` and `r.URL.Path`. Respond via
`jsonResponse(w, ...)` with a fixture body matching the spec.

```go
func TestAccountRetrieve(t *testing.T) {
    c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet || r.URL.Path != "/v1/account" {
            t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
        }
        jsonResponse(w, `{"id":"acct_1","name":"Acme","timezone":"UTC","created_at":"2026-01-01T00:00:00Z"}`)
    })
    got, err := c.Account.Retrieve(context.Background())
    if err != nil {
        t.Fatalf("Retrieve: %v", err)
    }
    if got.ID != "acct_1" || got.Name != "Acme" {
        t.Fatalf("unexpected payload: %+v", got)
    }
}
```

**Error path.** Use `unauthorizedHandler` from `testhelpers_test.go` and
assert with `expectAuthError(t, err)`.

```go
func TestAccountRetrieveUnauthorized(t *testing.T) {
    c, _ := newTestServer(t, unauthorizedHandler)
    _, err := c.Account.Retrieve(context.Background())
    expectAuthError(t, err)
}
```

### Decode test for required fields

For every resource file that defines a typed response, include one test where
the server returns an empty `{}` and verify the typed struct decodes with
zero-values *without* panicking — Go's JSON decoder is lenient by design, so
the test exists to lock in that behavior rather than to assert a specific
error.

## 5. After regeneration

The controller (`regenerate.sh` via `tools/regenerate-go.sh`) will run:

```bash
gofmt -l .                       # must produce no output (zero unformatted files)
go vet ./...
go build ./...
go test ./...
```

All four must pass. If any fail, the regeneration is rejected; fix the
underlying `CLAUDE.md` bug rather than the generated output.

## 6. Version bump

After a successful regen with API surface changes, bump the `Version`
constant in `version.go`:

- Patch (`1.0.0` → `1.0.1`) — bug fixes / docstring tweaks only.
- Minor (`1.0.0` → `1.1.0`) — backward-compatible additions (new endpoints, new optional fields).
- Major (`1.0.0` → `2.0.0`) — breaking changes (renamed/removed fields, changed types).

Tag the commit `v<version>` and push the tag — the Go module proxy
(`pkg.go.dev`) auto-discovers the new version.
