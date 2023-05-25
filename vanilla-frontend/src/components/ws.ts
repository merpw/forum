import { WebSocketResponse } from "./types"

const WS_URL = "ws://localhost:6969"
export let ws: WebSocket;

export const wsHandler = async () => {
  ws = new WebSocket(WS_URL)

  ws.onmessage = (event) => {
    try {
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

      // const { handler } = chatHandlers.find(({ regex }) => data.item.url.match(regex)) ?? {}

      // if (!handler) {
      //   console.error("unhandled ws response", data)
      //   return
      // }
      // TODO: Figure out implementation of commented code above
    } catch (e) {
      console.error("ws error", e)
    }
  }
  ws.onopen = () => {
    console.log("ws connected")
    const token = document.cookie.match(/forum-token=(.*?)(;|$)/)?.[1]
    if (!token) {
      console.error("Not logged in")
      return
    }
    ws.send(JSON.stringify({ type: "handshake", item: { token } }))
  }
  ws.onclose = () => {
    console.log("ws disconnected")
  }
  

}


