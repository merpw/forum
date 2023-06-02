import { WSGetResponse, WSPostResponse, WebSocketResponse } from "./types"

const WS_URL = "ws://localhost:6969"
export let ws: WebSocket
let opened = false

export const wsHandler = () => {
  ws = new WebSocket(WS_URL)

  ws.onmessage = (event) => {
    try {
      // const request = JSON.parse(event.data)
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

      if (data.type === "get") {
        getHandler(data)
      }

      if (data.type === "post") {
        postHandler(data)
      }
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

    if (!opened) {
      opened = true
      ws.send(JSON.stringify({ type: "handshake", item: { token } }))
    }
  }

  ws.onclose = () => {
    console.log("ws disconnected")
  }
}

const postHandler = (resp: WSPostResponse<WebSocketResponse<never>>) => {
  const respData = resp as Object
  console.log(respData)

}

const getHandler = (resp: WSGetResponse<WebSocketResponse<never>>) => {
  const respData = resp as Object
  console.log(respData)

}
