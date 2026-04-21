# strm-router

**strm-router** is a lightweight HTTP router that forwards incoming requests to different backend services based on the `Authorization` header. It is designed for simple service routing in cloud-based environments.

---

## Features

* Routes requests based on `Authorization` header values
* Simple JSON-based configuration
* Minimal dependency setup
* Fast Go-based HTTP proxying

---

## Configuration

Create a `config.json` file in the project root:

```json
{
  "serviceA": "http://localhost:8081",
  "serviceB": "http://localhost:8082"
}
```

---

## Build Instructions

Make sure you have Go installed (1.18+ recommended).

```bash
go build main.go
```

This will produce an executable binary.

---

## Run

```bash
./main
```

By default, the server starts on the configured port (check source code if you need to change it).

---

## Request Flow

1. Client sends HTTP request with `Authorization` header
2. `strm-router` reads the header value
3. It matches the value against `config.json`
4. Request is proxied to the mapped service URL
5. Response is returned to the client

---

## Example Request

```bash
aws s3 ls
```
```bash
strm bkt ls
```

---

## License

[MIT](./LICENSE)
