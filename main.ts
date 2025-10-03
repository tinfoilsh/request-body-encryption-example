import { SecureClient } from "tinfoil";

declare global {
  interface Window {
    tinfoil?: {
      client: SecureClient;
      ready: Promise<void>;
    };
  }
}

// Use inference host here
const client = new SecureClient({
  baseURL: "http://localhost:8080/",
  hpkeKeyURL: "https://ehbp.inf6.tinfoil.sh/v1/",
  configRepo: "tinfoilsh/confidential-inference-proxy-hpke",
});

const ready = client.ready();

// expose to window for HTML to use
window.tinfoil = { client, ready };