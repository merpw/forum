import { chatList, currentChat, messages, renderChatList } from "./chat.js"
import { WSGetResponse, WSPostResponse, WebSocketResponse, Message } from "../types"
import { iterator } from "./utils.js"
import { Auth, userInfo } from "./auth.js"

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
        return
      }

      if (data.type === "post") {
        postHandler(data)
        return
      }

    } catch (e) {
      console.error("ws error", e)
      ws.close()
      setTimeout(() => {
        Auth(false)
      }, 25)
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

const postHandler = async (resp: WSPostResponse<never>) => {
  const data = resp.item.data
  const url = resp.item.url

  if (url.match(/^\/chat\/\d{1,}\/message$/)) { // /chat/{id}/message
    getMessage(data)
  }
}

const getHandler = (resp: WSGetResponse<never>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url

  /* /chat/{id} */
  if (url.match(/^\/chat\/\d{1,}$/)) {
    setTimeout(() => {
      for (const user of Object.values(chatList.Users)) {
        if (user.Id == data.companionId) {
          console.log("chat", data)
          Object.assign(user, {LastMessageId: data.lastMessageId})
          break
        }
      }

      if (!chatList.Ids.has(data.companionId)) {
        chatList.Ids.set(data.companionId, data.id)
        chatList.Users.forEach((user) => {
          if (user.Id == data.companionId) {
            user.ChatId = data.id
            return
          }
        })
      }
    }, 20)
  }

  /* /chat/{id}/messages */
  if (url.match(/^\/chat\/\d{1,}\/messages$/)) { 
    getMessageList(data as number[])
    return
  }

  /* /message/{id} */
  if (url.match(/^\/message\/\d{1,}$/)) {
    if (data.authorId != -1) {
      chatList.Users.forEach((user) => {
        if (user.ChatId == data.chatId) {
          user.LastMessageId = data.id
          return
        }
      }) 
    }

    setTimeout(() => {
      messages.current = data as Message
      if (data.chatId === currentChat.chatId) {      
        messages.list.unshift(data as Message)
      } else {
        updateUnreadMessages()
      }
    }, 20)
    return
  }
 /* /users/online */
  if (url.match(/^\/users\/online$/)) {
    updateOnlineUsers(data as number[])
    return
  }

  if (url.match(/^\/chat\/all$/)) {
    getChatIds(resp)
    return
  }
}

function updateOnlineUsers(users: number[]) {
    for (const user of Object.values(chatList.Users)) {
      if (users.includes(user.Id)) {
        user.Online = true
      } else {
        user.Online = false
      }
    } 
  renderChatList()
}

function updateUnreadMessages() {
  for (const [userId, chatId] of chatList.Ids) {
    if (chatId === messages.current.chatId) {
      for (const user of Object.values(chatList.Users)) {
        if (user.Id === userId) {
          user.UnreadMsg = true
          break
        }
      }
      break
    }
  }
  renderChatList()
}

const getChatIds = (resp: iterator) => {
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

const getMessage = (id: number) => {
  sendWsObject({
    type: "get",
    item: {
      url: `/message/${id}`,
    },
  })
  setTimeout(() => {
    const chat = document.getElementById(`Chat${currentChat.chatId}`) as HTMLDivElement
        if (!chat) {
          } else {
      chat.dispatchEvent(messageEvent)
    }
  }, 80)
}

const getMessageList = (ids: number[]) => {
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
