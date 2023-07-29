import { client } from "../../main.js"
import {
  WSGetResponse,
  WSPostResponse,
  WebSocketResponse,
  Message,
  User,
} from "../../types"
import { Auth, state } from "../authorization/auth.js"
import { iterator } from "../utils.js"
import { sendWsObject } from "./helpers/sendobject.js"
import {
  renderChatList,
  renderNewMessages,
  startChat,
} from "./helpers/events.js"
import { Chat } from "./chat.js"
import { getUserById } from "../../api/get.js"

export let ws: WebSocket
const WSUrl = "ws://localhost/ws"

export const wsHandler = () => {
  let retryTimeout = 0
  ws = new WebSocket(WSUrl)
  ws.onmessage = (event: MessageEvent): void => {
    try {
      const data = JSON.parse(event.data) as WebSocketResponse<never>

      if (data.type === "handshake") {
        console.log("ws handshake success")
        window.dispatchEvent(startChat)
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
      Auth(false)
    }
  }

  ws.onopen = () => {
    console.log("ws connected")
    const token = document.cookie.match(/forum-token=(.*?)(;|$)/)?.[1]
    if (!token) {
      console.error("Not logged in")
      return
    }
    retryTimeout = 0
    ws.send(JSON.stringify({ type: "handshake", item: { token } }))
  }

  ws.onclose = () => {
    console.log("ws disconnected")
    if (document.cookie.match(/forum-token=(.*?)(;|$)/)?.[1]) {
      setTimeout(() => {
        retryTimeout += 1000
        wsHandler()
      }, retryTimeout)
    }
  }
}

const postHandler = (resp: WSPostResponse<never>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url

  // Handle all messages sent.
  if (url.match(/^\/chat\/\d+\/message$/)) {
    const chatId = +url.split("/")[2]
    if (!client.chatIds?.includes(chatId)) return

    sendWsObject({
      type: "get",
      item: {
        url: `/message/${data}`,
      },
    })

    const chat = client.chats.get(chatId) as Chat
    const messageId = parseInt(`${data}`)
    chat.lastMessageId = messageId
    if (!chat.messageIds.includes(messageId)) {
      chat.messageIds.unshift(messageId)
    }

    // This is just for notifications. Will add if chat is not opened or if you have scrolled too far up.
    if (chat.chatId === client.activeChat?.chatId) {
      if (chat.chatMessages.scrollTop <= -50) {
        chat.unreadMessages += 1
      }

      chat.chatMessages.dispatchEvent(renderNewMessages)
    } else {
      chat.unreadMessages += 1
    }

    window.dispatchEvent(renderChatList)
  }

  // CREATE a chat if not existing in DB.
  // ADD chatId to list of chatIds
  // also ADD to client object as a Chat object.
  if (url.match(/^\/chat\/create$/)) {
    const chatId = data.chatId as number

    if (!client.chatIds?.includes(chatId)) {
      client.chatIds?.unshift(chatId)
      sendWsObject({
        type: "get",
        item: {
          url: `/chat/${chatId}`,
        },
      })
     
    }
  }
}

const getHandler = async (resp: WSGetResponse<never>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url

  // ASSIGN chats to all IDs from the response array.
  // Example:
  // [1, 2, 3] => (1 => Chat, 2 => Chat, 3 => Chat)
  if (url.match(/^\/chat\/all$/)) {
    const chatIds = data as number[]
    client.chatIds = chatIds

    for (const id of chatIds) {
      sendWsObject({
        type: "get",
        item: {
          url: `/chat/${id}`,
        },
      })
    }

    window.dispatchEvent(renderChatList)
  }

  // ADD new Chat if not in memory, UPDATE lastMessageId if already in memory
  // Chat will be added as following format: Name, userId, chatId, lastMessageId
  if (url.match(/^\/chat\/\d+$/)) {
    client.userChats.set(data.companionId, data.id)

    if (!client.chats.has(data.id as number)) {
      const user = state.users.get(data.companionId as number) as User
      client.chats.set(
        data.id,
        new Chat(
          user.name,
          data.companionId,
          data.id,
          data.lastMessageId,
          client.onlineUsers.includes(user.id)
        )
      )
    } else {
      const chat = client.chats.get(data.id) as Chat
      chat.lastMessageId = data.lastMessageId
    }

    sendWsObject({
      type: "get",
      item: {
        url: `/chat/${data.id}/messages`,
      },
    })
  }

  // ADD messageIds to Chat
  if (url.match(/^\/chat\/\d+\/messages$/)) {
    const id = +url.split("/")[2]
    const chat = client.chats.get(id)
    if (!chat) return

    const messageIds = data.reverse() as number[]
    chat.messageIds = messageIds
  }

  // SET client messages with this format: MessageId => Message
  if (url.match(/^\/message\/\d+$/)) {
    const id = +url.split("/")[2]
    client.messages.set(id, data as Message)
  }

  // UPDATE list of online users in memory. Dispatch event to render chat list.
  if (url.match(/^\/users\/online$/)) {
    const onlineIds = data as number[]
    const filteredIds = onlineIds.filter(
      (id) => id !== state.me.id && id !== -1
    )

    for (const id of filteredIds) {
      if (!state.users.has(id)) {
        state.users.set(id, await getUserById(id))
      }
    }

    client.onlineUsers = filteredIds

    sendWsObject({
      type: "get",
      item: {
        url: "/chat/all",
      },
    })
  }
}
