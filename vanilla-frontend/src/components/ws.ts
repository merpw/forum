import { chatList, messages } from "./chat.js"
import { WSGetResponse, WSPostResponse, WebSocketResponse } from "../types"

const WS_URL = "ws://localhost:6969"
export let ws: WebSocket


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

      if (data.type === "get") {
        getHandler(data)
      }

      if (data.type === "post") {
        console.log(data)
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
    ws.send(JSON.stringify({ type: "handshake", item: { token } }))
  }

  ws.onclose = () => {
    console.log("ws disconnected")
  }
}

const postHandler = (resp: WSPostResponse<WebSocketResponse<never>>) => {
 console.log(resp) 
}

const getHandler = async (resp: WSGetResponse<any>) => {
  const chatIds = new RegExp(/^\/chat\/\d{1,}$/)
  const messageList = new RegExp(/^\/chat\/\d{1,}\/messages$/)
  const message = new RegExp(/^\/message\/\d{1,}$/)

  if (resp.item.url.match(chatIds)) {
    chatList.Ids.set(resp.item.data.userId, resp.item.data.id,)
    return
  }

  if (resp.item.url.match(messageList)) {
    await getMessageList(resp.item.data)
  }

  if (resp.item.url.match(message)) {
    messages.list.push(resp.item.data) 
  }

  if (resp.item.url === "/chat/all") {
    getChatIds(resp) 
    return
  }
}

const getChatIds = async (resp: WSGetResponse<any>) => {

  for (const user of resp.item.data) {
    console.log(user)
    sendWsObject(
      {
        type: "get",
        item: {
          url: `/chat/${user}`
        }
      }
    )
  }
  return
}

const getMessageList = async (ids: any[]) => {
  for (const id of ids) {
    console.log(id)
    sendWsObject(
      {
        type: "get",
        item: {
          url: `/message/${id}`
        }
      }
    )
  }
  return
}

export async function sendWsObject(obj: unknown) {
  if (ws.readyState === 1) {
    ws.send(JSON.stringify(obj))
  }
}
