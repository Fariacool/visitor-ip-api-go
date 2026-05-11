# visitor-ip-api-go

A small Go API that returns the best-known client public IP address.

## Features

- `GET /ip` returns only the resolved client IP.
- Client IP selection uses a simple best-effort chain: `CF-Connecting-IP` / `CF-Connecting-IPv6`, then `X-Forwarded-For`, then `X-Real-IP`, and finally `RemoteAddr`.
- Generated `OpenAPI 3.1` and `Swagger UI` are exposed by Huma at `/openapi.json`, `/openapi.yaml`, and `/docs`.
- Built-in CORS for browser callers.

## Response Example

```json
{
  "ip": "2606:4700:4700::1111"
}
```

## Environment Variables

- `LISTEN_ADDR`
  - Default: `:8466`

## Run

```bash
go run .
```

Then open:

- `http://127.0.0.1:8466/docs`
- `http://127.0.0.1:8466/openapi.json`

## Docker

```bash
docker run --rm -p 8466:8466 ghcr.io/fariacool/visitor-ip-api-go:latest
```

With Compose:

```bash
docker compose up -d
```

The repository includes:

- `Dockerfile` for a multi-stage container build
- `docker-compose.yml` for local container runs
- `.github/workflows/release.yml` to publish `linux/amd64` and `linux/arm64` images to `ghcr.io`

## Cloudflare Notes

- If Cloudflare `Pseudo IPv4` is set to `Overwrite Headers`, Cloudflare replaces `CF-Connecting-IP` with a pseudo IPv4 value and preserves the real IPv6 in `CF-Connecting-IPv6`. This service returns the IPv6 address in that case.
- This service does not treat proxy headers as a security boundary. It simply returns the best-known client IP from common proxy headers before falling back to `RemoteAddr`.
