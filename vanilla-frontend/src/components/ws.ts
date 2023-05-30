import { WebSocketResponse } from "./types"

const WS_URL = "ws://localhost:6969"
export let ws: WebSocket
let opened = false

export const wsHandler = () => {
  ws = new WebSocket(WS_URL)
  ws.onopen = () => {
    console.log("ws connected")
    const token = document.cookie.match(/forum-token=(.*?)(;|$)/)?.[1]
    if (!token) {
      console.error("Not logged in")
      return
    }
    if (!opened){
      opened = true
      ws.send(JSON.stringify({ type: "handshake", item: { token } }))
    }
  }

  ws.onmessage = (event) => {
    try {
      console.log("Event:", event)
      const data = JSON.parse(event.data) as WebSocketResponse<never>

      if (data.type === "handshake") {
        console.log("ws handshake success")
        return
      }

      if (data.type === "error") {
        console.error("ws error:", data.item.message)
        return
      }

      if (!data.item?.url) {
        console.error("invalid ws response:", data)
        return
      }

      // TODO: Figure out implementation of commented code above
    } catch (e) {
      console.error("ws error", e)
    }
  }

  ws.onclose = () => {
    console.log("ws disconnected")
  }
}
