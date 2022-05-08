export const ENDPOINT = import.meta.env.PROD
  ? window.location.origin
  : "http://localhost:8080";

export const HTTP_ENDPOINT = ENDPOINT + "/api/v1";

export const WS_ENDPOINT =
  (ENDPOINT.startsWith("https://") ? "wss://" : "ws://") +
  ENDPOINT.split("://")[1] +
  "/ws";

console.log(ENDPOINT, WS_ENDPOINT);
