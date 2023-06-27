import { chatList, currentChat, messages } from "./chat.js"
import { WSGetResponse, WSPostResponse, WebSocketResponse, Message } from "../types"
import { iterator } from "./utils.js"
import { Auth } from "./auth.js"

export const messageEvent = new Event("messageEvent")
const WS_URL = "ws://localhost:8081/ws"
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

const postHandler = async (resp: WSPostResponse<object | any>) => {
  const data = resp.item.data
  const url = resp.item.url
  // const createChat = new RegExp(/^\/chat\/create$/)
  if (url.match(/^\/chat\/\d{1,}\/message$/)) { // /chat/{id}/message
    getMessage(data)
  }
}

const getHandler = async (resp: WSGetResponse<object>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url

  if (url.match(/^\/chat\/\d{1,}$/)) {
    if (!chatList.Ids.has(data.companionId)) {
      return chatList.Ids.set(data.companionId, data.id)
    }
  }

  if (url.match(/^\/chat\/\d{1,}\/messages$/)) { // chat/{id}/messages
    getMessageList(data as number[])
    return
  }

  if (url.match(/^\/message\/\d{1,}$/)) { // /message/{id}
    if (data.chatId === currentChat.chatId) {
      messages.list.unshift(data as Message)
    }
    return
  }

  if (url.match(/^\/chat\/all$/)) {
    getChatIds(resp)
    return
  }
}

const getChatIds = async (resp: iterator) => {
  for (const chat of resp.item.data) {
    if (!chatList.Ids.has(chat)) {
      sendWsObject({
        type: "get",
        item: {
          url: `/chat/${chat}`,
        },
      })
    }
  }
  return
}

const getMessage = async (id: number) => {
  sendWsObject({
    type: "get",
    item: {
      url: `/message/${id}`,
    },
  })
  setTimeout(() => {
    document.getElementById(`Chat${currentChat.chatId}`)?.dispatchEvent(messageEvent)
  }, 150)

  return
}

const getMessageList = async (ids: number[]): Promise<void> => {
  for (const id of ids) {
      sendWsObject({
      type: "get",
      item: {
        url: `/message/${id}`,
      },
    })
  }
  return
}

export function sendWsObject(obj: object) {
    if (ws.readyState === 1) {
      ws.send(JSON.stringify(obj))
    }
  return
}
