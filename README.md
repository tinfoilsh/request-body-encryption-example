# OpenAI Example (HPKE)

Simple browser chat that encrypts request bodies end-to-end with HPKE before reaching the backend.

## Prerequisites

- Go 1.24+
- Node modules already installed under `../tinfoil-node` (provides the EHBP bundle)
- `TINFOIL_API_KEY` (or `OPENAI_API_KEY`) for the upstream enclave/proxy

## Run the backend

```bash
cd openai-example
TINFOIL_API_KEY=your_key go run .
```

This launches the HPKE-aware proxy on `http://localhost:8080`, publishing its key at `/.well-known/hpke-keys` and forwarding to the enclave (override with `ENCLAVE_CHAT_URL` if needed).

## Serve the frontend

From the repository root:

```bash
python -m http.server 8000
```

Then open `http://localhost:8000/openai-example/index.html` in your browser. Every chat message will be HPKE-encrypted before leaving the page.
