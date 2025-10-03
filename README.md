# Simple AI Chat with Tinfoil Security

This is a simple chat interface that uses the Tinfoil SecureClient to communicate securely with an AI backend.

## Prerequisites

- Go installed (for running the backend server)
- A Tinfoil API key (set as `TINFOIL_API_KEY` in your environment)

## Quick Start

1. Set your Tinfoil API key as an environment variable:
   ```bash
   export TINFOIL_API_KEY=your_api_key_here
   ```

2. Start the backend server:
   ```bash
   go run main.go
   ```

3. Serve the frontend files:
   ```bash
   # Using Python 3
   python3 -m http.server 8000
   
   # Or using Node.js (if you have http-server installed)
   # npx http-server -p 8000
   ```

4. Open your browser and navigate to http://localhost:8000

## How It Works

The application uses the Tinfoil SecureClient to establish a secure connection with the backend server. All communication between the frontend and backend is encrypted using HPKE (Hybrid Public Key Encryption).

The backend server acts as a proxy between the frontend and the Tinfoil inference service, forwarding requests while adding the necessary authentication headers.

## Customization

You can modify the following parameters in `index.html`:

- Model selection (currently set to 'gpt-oss-120')
- Backend URL (currently set to 'http://localhost:8080')
- HPKE key URL and config repository

## Troubleshooting

If you encounter issues:

1. Check that your TINFOIL_API_KEY is correctly set
2. Verify that the backend server is running on port 8080
3. Check the browser console for any JavaScript errors
4. Ensure that your API key has the necessary permissions

For more information about the Tinfoil SecureClient, visit [Tinfoil Documentation](https://docs.tinfoil.dev).