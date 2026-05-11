# OpenResty Configs

This directory contains ready-to-use OpenResty / nginx config examples for `visitor-ip-api-go`.

## Files

- `visitor-ip-api-go-standalone-server.conf`
  - Use when `visitor-ip-api-go` gets its own dedicated `server {}` block.
  - Place it in a directory that is included from the `http {}` context, for example `conf.d/`.

## Upstream Address

The config proxies to:

```nginx
http://127.0.0.1:8466
```

Change that if your Go service listens somewhere else.

## Headers Passed Upstream

The config forwards:

- `Host`
- `X-Real-IP`
- `X-Forwarded-For`
- `X-Forwarded-Proto`
- `X-Forwarded-Host`
- `X-Forwarded-Port`
- Cloudflare headers used by the app:
  - `CF-Connecting-IP`
  - `CF-Connecting-IPv6`

## Reload

After copying the config:

```bash
openresty -t
```

Then reload OpenResty using your platform's service manager or process control command.

If your package uses `nginx` instead of `openresty`, use:

```bash
nginx -t
```

Then reload nginx using your platform's service manager or process control command.

## CORS

This config file already includes enhanced CORS handling at the OpenResty layer:

- hides upstream `Access-Control-Allow-*` headers to avoid duplicates
- returns `204` for `OPTIONS` without hitting the Go app
- adds CORS headers with `always`, so proxy-generated `4xx` / `5xx` responses are also cross-origin readable

In other words, OpenResty is the source of truth for CORS in this config.
