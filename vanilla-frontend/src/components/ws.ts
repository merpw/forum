import { chatList, messages } from "./chat.js"
import { WSGetResponse, WSPostResponse, WebSocketResponse } from "../types"
import { iterator } from "./utils.js"
import { Auth } from "./auth.js"

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
      ws.close()
      setTimeout(() => {
        Auth(false)
      }, 50)
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

const postHandler = (resp: WSPostResponse<any>) => {
  const data: iterator = resp.item.data
  const createChat = new RegExp(/^\/chat\/create$/)
  console.log(data)
}

const getHandler = async (resp: WSGetResponse< object | any >) => {
  const data: iterator = resp.item.data
  const url = resp.item.url
  const chatIds = new RegExp(/^\/chat\/\d{1,}$/)
  const messageList = new RegExp(/^\/chat\/\d{1,}\/messages$/)
  const message = new RegExp(/^\/message\/\d{1,}$/)

  if (url.match(chatIds)) {
    if (!chatList.Ids.has(data.userId)){
    return chatList.Ids.set(data.userId, data.id)
    }
    
  }

  if (url.match(messageList)) {
    return await getMessageList(resp.item.data)
  }

  if (url.match(message)) {
    return messages.list.push(resp.item.data) 
  }

  if (url === "/chat/all") {
    return await getChatIds(resp) 
  }
}

const getChatIds = async (resp: iterator) => {
  for (const user of resp.item.data) {
    if (!chatList.Ids.has(user)){
    sendWsObject(
      {
        type: "get",
        item: {
          url: `/chat/${user}`
        }
      }
    )

    }
  }
  return
}

const getMessageList = async (ids: number[]) => {
  for (const id of ids) {
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
  return
}
