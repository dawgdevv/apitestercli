# Writing Test Files for Probe

This guide covers everything you need to write YAML test files for Probe.

---

## File Structure

Every test file has two top-level keys:

```yaml
env:
  # variables go here

tests:
  # test cases go here
```

---

## `env` — Environment Variables

The `env` block defines reusable variables for your test suite. The only **required** variable is `base_url`.

```yaml
env:
  base_url: https://api.example.com
  api_key: my-secret-key
  user_id: "42"
```

- Keys are arbitrary — name them whatever makes sense
- Values must be strings
- You can reference variables anywhere using `{{variable_name}}` syntax (paths, headers, body)

---

## `tests` — Test Cases

Each test case is an item in the `tests` array with three fields:

| Field | Required | Description |
|---|---|---|
| `name` | Yes | Human-readable test name (shown in output) |
| `request` | Yes | The HTTP request to make |
| `expect` | Yes | What the response should look like |

Minimal example:

```yaml
tests:
  - name: Health check
    request:
      method: GET
      path: /health
    expect:
      status: 200
```

---

## `request` — Defining the HTTP Request

| Field | Required | Type | Description |
|---|---|---|---|
| `method` | Yes | String | HTTP method: `GET`, `POST`, `PUT`, `DELETE` |
| `path` | Yes | String | URL path (appended to `base_url`) |
| `headers` | No | Map | Key-value pairs for request headers |
| `body` | No | Map | JSON request body |

### GET Request

```yaml
request:
  method: GET
  path: /users/1
```

### POST Request with Headers and Body

```yaml
request:
  method: POST
  path: /posts
  headers:
    Content-Type: application/json
    Authorization: Bearer {{api_key}}
  body:
    title: "Hello World"
    body: "Post content"
    userId: "1"
```

### PUT Request

```yaml
request:
  method: PUT
  path: /posts/1
  headers:
    Content-Type: application/json
  body:
    id: "1"
    title: "Updated title"
    body: "Updated body"
    userId: "1"
```

### DELETE Request

```yaml
request:
  method: DELETE
  path: /posts/1
```

---

## `expect` — Assertions

| Field | Required | Type | Description |
|---|---|---|---|
| `status` | Yes | Integer | Expected HTTP status code |
| `json` | No | Map | JSON field assertions on the response body |

### Status Code Only

```yaml
expect:
  status: 200
```

```yaml
expect:
  status: 201
```

```yaml
expect:
  status: 404
```

### JSON Field Assertions

Assert exact values on fields in the JSON response:

```yaml
expect:
  status: 200
  json:
    id: 1
    username: "Bret"
```

---

## JSON Assertion Types

### Exact Match

Check that a field equals a specific value:

```yaml
json:
  id: 1
  title: "Hello World"
```

### Nested Fields (Dot Notation)

Access deeply nested JSON fields using dots:

```yaml
# Response: { "address": { "city": "Gwenborough", "zipcode": "92998-3874" } }
json:
  "address.city": "Gwenborough"
  "address.zipcode": "92998-3874"
  "company.name": "Romaguera-Crona"
```

> **Note:** Wrap dot-notation keys in quotes so YAML parses them as a single string.

### Array Length

Check the length of a JSON array response using `$.length`:

```yaml
# Assert the response array has more than 50 items
json:
  "$.length": ">50"
```

```yaml
# Assert the response array has more than 5 items
json:
  "$.length": ">5"
```

### Numeric Comparisons

Use `>` and `<` for numeric comparisons:

```yaml
json:
  "$.length": ">100"
```

---

## Variable Substitution

Use `{{variable_name}}` anywhere in paths, headers, or body values. Variables are defined in the `env` block.

```yaml
env:
  base_url: https://jsonplaceholder.typicode.com
  user_id: "3"
  post_id: "25"
  content_type: application/json
  post_title: "My Post"

tests:
  - name: Get user by ID
    request:
      method: GET
      path: /users/{{user_id}}
    expect:
      status: 200

  - name: Create post with variables
    request:
      method: POST
      path: /posts
      headers:
        Content-Type: "{{content_type}}"
      body:
        title: "{{post_title}}"
        userId: "{{user_id}}"
    expect:
      status: 201
```

Variables can reference other variables:

```yaml
env:
  post_title: "Substituted Title"
  post_body: "This body references {{post_title}}"
```

---

## Complete Examples

### Basic CRUD Suite

```yaml
env:
  base_url: https://jsonplaceholder.typicode.com

tests:
  - name: List all posts
    request:
      method: GET
      path: /posts
    expect:
      status: 200

  - name: Get single post
    request:
      method: GET
      path: /posts/1
    expect:
      status: 200
      json:
        id: 1
        userId: 1

  - name: Create a post
    request:
      method: POST
      path: /posts
      headers:
        Content-Type: application/json
      body:
        title: "Probe test post"
        body: "Created by Probe"
        userId: "1"
    expect:
      status: 201
      json:
        title: "Probe test post"

  - name: Update a post
    request:
      method: PUT
      path: /posts/1
      headers:
        Content-Type: application/json
      body:
        id: "1"
        title: "Updated title"
        body: "Updated body"
        userId: "1"
    expect:
      status: 200
      json:
        title: "Updated title"

  - name: Delete a post
    request:
      method: DELETE
      path: /posts/1
    expect:
      status: 200
```

### Assertion-Heavy Suite

```yaml
env:
  base_url: https://jsonplaceholder.typicode.com

tests:
  - name: Exact field match
    request:
      method: GET
      path: /posts/1
    expect:
      status: 200
      json:
        id: 1
        userId: 1

  - name: Nested field match
    request:
      method: GET
      path: /users/1
    expect:
      status: 200
      json:
        username: "Bret"
        "address.city": "Gwenborough"
        "company.name": "Romaguera-Crona"

  - name: Array length check
    request:
      method: GET
      path: /posts
    expect:
      status: 200
      json:
        "$.length": ">50"
```

### Error Handling Suite

```yaml
env:
  base_url: https://jsonplaceholder.typicode.com

tests:
  - name: 404 — Non-existent resource
    request:
      method: GET
      path: /posts/99999
    expect:
      status: 404

  - name: 404 — Invalid endpoint
    request:
      method: GET
      path: /nonexistent
    expect:
      status: 404
```

---

## Running Your Tests

```bash
# Run from CLI
probe run tests.yaml

# Run from web dashboard
probe serve
```

---

## Quick Reference

```yaml
env:
  base_url: https://api.example.com     # Required
  any_variable: "value"                  # Optional, reusable

tests:
  - name: "Test name"                    # Required
    request:
      method: GET                        # GET | POST | PUT | DELETE
      path: /endpoint/{{any_variable}}   # Supports {{var}} substitution
      headers:                           # Optional
        Content-Type: application/json
      body:                              # Optional (POST/PUT)
        key: "value"
    expect:
      status: 200                        # Required — HTTP status code
      json:                              # Optional — response assertions
        field: "exact value"             # Exact match
        "nested.field": "value"          # Dot notation
        "$.length": ">10"               # Array length comparison
```
