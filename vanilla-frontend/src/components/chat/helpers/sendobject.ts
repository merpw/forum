import { ws } from "../ws.js"

// Sends a WebSocket object as a message to the WebSocket server
export function sendWsObject(obj: object): void {
  if (ws.readyState === 1) {
    ws.send(JSON.stringify(obj))
  }
}
