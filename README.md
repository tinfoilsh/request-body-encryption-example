# Request Body Encryption Example

A small example that demonstrates how to use the Tinfoil Request Body Encryption feature from a simple browser-based chat. It
contains two pieces:

- `main.go`: a Go proxy that adds your TINFOIL_API_KEY and forwards chat completions to
the Tinfoil enclaves for inference.
- `main.ts`: a few lines of TypeScript that instantiates the Tinfoil `SecureClient`, sends
  `/v1/chat/completions`, and streams the response into the page.

In this example, the Go proxy only handles `/v1/chat/completions`. All other requests are handled directly by the Tinfoil backend and enclaves.

The Tinfoil enclave that all chat completions requests must be forwarded to is currently accessible at `https://ehbp.inf6.tinfoil.sh/v1/chat/completions`.

## Prerequisites

- `npm install` already run in this directory
- `TINFOIL_API_KEY` exported in your shell

## Running the demo

```bash
# Terminal 1 – start the proxy on http://localhost:8080
export TINFOIL_API_KEY=tk-...
go run main.go

# Terminal 2 – serve the static files with Vite
npx vite
```

Open the printed Vite URL (typically http://localhost:5173), type a message, and
watch the assistant stream its reply. The demo intentionally keeps the UI and
error handling minimal so it is easy to read and adapt.

### Tweaks

- Change the `baseURL` or model in `main.ts` if you want to point at a different
  proxy server or model. Defaults are `http://localhost:8080` and `gpt-oss-120b`.
