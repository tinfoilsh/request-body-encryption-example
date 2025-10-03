# Tinfoil Chat Demo

A tiny example that shows how to use the Tinfoil `SecureClient` from a browser. It
contains two pieces:

- `main.go`: a Go proxy that adds your API key and forwards chat completions to
the hosted HPKE service while preserving the EHBP headers.
- `main.ts`: a few lines of TypeScript that instantiates `SecureClient`, sends
  `/v1/chat/completions`, and streams the response into the page.

## Prerequisites

- `npm install` already run in this directory
- `TINFOIL_API_KEY` exported in your shell

## Running the demo

```bash
# Terminal 1 – start the proxy on http://localhost:8080
export TINFOIL_API_KEY=sk-...
go run main.go

# Terminal 2 – serve the static files with Vite
npx vite
```

Open the printed Vite URL (typically http://localhost:5173), type a message, and
watch the assistant stream its reply. The demo intentionally keeps the UI and
error handling minimal so it is easy to read and adapt.

### Tweaks

- Change the `baseURL` or model in `main.ts` if you want to point at a different
  proxy or model.
- The proxy currently targets `https://ehbp.inf6.tinfoil.sh`; edit `main.go` if
  you are testing another environment.
