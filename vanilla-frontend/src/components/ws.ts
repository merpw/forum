import { chatList, currentChat, messages } from "./chat.js"
import { WSGetResponse, WSPostResponse, WebSocketResponse } from "../types"
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

const postHandler = async (resp: WSPostResponse<object | any>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url
  // const createChat = new RegExp(/^\/chat\/create$/)
  const sendMessage = new RegExp(/^\/chat\/\d{1,}\/message$/)
  if (url.match(sendMessage)) {
    getMessage(data.messageId)
  }
}

const getHandler = async (resp: WSGetResponse<object | any>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url
  const chatIds = new RegExp(/^\/chat\/\d{1,}$/)
  const messageList = new RegExp(/^\/chat\/\d{1,}\/messages$/)
  const message = new RegExp(/^\/message\/\d{1,}$/)

  if (url.match(chatIds)) {
    if (!chatList.Ids.has(data.companionId)) {
      return chatList.Ids.set(data.companionId, data.id)
    }
  }

  if (url.match(messageList)) {
    return await getMessageList(resp.item.data)
  }

  if (url.match(message)) {
    if (resp.item.data.chatId === currentChat) {
      messages.list.push(resp.item.data)
    }
    return
  }

  if (url === "/chat/all") {
    return await getChatIds(resp)
  }
}

const getChatIds = async (resp: iterator): Promise<void> => {
  for (const chat of resp.item.data) {
    if (!chatList.Ids.has(chat)) {
      await sendWsObject({
        type: "get",
        item: {
          url: `/chat/${chat}`,
        },
      })
    }
  }
  return
}

const getMessage = async (id: number): Promise<void> => {
  await sendWsObject({
    type: "get",
    item: {
      url: `/message/${id}`,
    },
  }).then(() => {
    document.getElementById(`Chat${currentChat}`)?.dispatchEvent(messageEvent)
  })
}

const getMessageList = async (ids: number[]): Promise<void> => {
  for (const id of ids) {
    await sendWsObject({
      type: "get",
      item: {
        url: `/message/${id}`,
      },
    })
  }
  return
}

export function sendWsObject(obj: object): Promise<void> {
  return new Promise((resolve, reject) => {
    if (ws.readyState === 1) {
      ws.send(JSON.stringify(obj))
      resolve()
    } else {
      reject(new Error("WebSocket connection is not open."))
    }
  })
}
