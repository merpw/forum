import { client } from "../../main.js"
import { WSGetResponse, WSPostResponse, WebSocketResponse, Message } from "../../types"

import { Auth } from "../authorization/auth.js"

import { iterator } from "../utils.js"
import { sendWsObject } from "./helpers/sendobject.js"
import { renderChatList, startChat, chatCreated } from "./helpers/events.js"
import { Chat } from "./client.js"
import { getUserById } from "../../api/get.js"

export let ws: WebSocket

export const wsHandler = () => {
  let retryTimeout = 0
  ws = new WebSocket(`${location.protocol.replace("http", "ws")}//${location.host}/ws`)
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
    if (document.cookie){
      setTimeout(() => {
        retryTimeout += 1000
        wsHandler() 
      }, retryTimeout)
    }
  }}

const postHandler = async (resp: WSPostResponse<never>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url

  if (url.match(/^\/chat\/\d+\/message$/)) { // /chat/{id}/message

  }

  // CREATES a chat if not existing in DB.
  // ADDS chatId to list of chatIds
  // also ADD to client object as a Chat object.
  if (url.match(/^\/chat\/create$/)) { // /chat/create
    console.log("chat/create", data)
    if(!client.chatIds?.includes(data.companionId as number)){
      client.chatIds?.unshift(data.chatId as number)
      sendWsObject({
        type: "get",
        item: {
          url: `/chat/${data.chatId as number}`
        }
      })
      window.dispatchEvent(chatCreated)
    }    
  }
}

const getHandler = async (resp: WSGetResponse<never>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url

  /* /chat/all */
  if (url.match(/^\/chat\/all$/)) {
    console.log("chat/all", data)
    const chatIds = data as number[]
    client.chatIds = chatIds
    for(const id of chatIds){
      sendWsObject({
        type: "get",
        item: {
          url: `/chat/${id}`
        }
      })
    }
    window.dispatchEvent(renderChatList)
  }

  // ADD new Chat if not in memory, UPDATE lastMessageId if already in memory
  // Chat will be added as following format: Name, userId, chatId, lastMessageId 
  if (url.match(/^\/chat\/\d+$/)) {
    console.log("chat/id", data)
    client.userChats.set(data.companionId, data.id)
    if(!client.chats.has(data.id as number)){
      // ADD chat 
      const user = await getUserById(data.companionId)
      client.chats.set(data.id, new Chat(user.Name, data.companionId, data.id, data.lastMessageId))
    } else {
      // UPDATE lastMessageId
      const chat = client.chats.get(data.id) as Chat
      chat.lastMessageId = data.lastMessageId
    }
  }

  // ADD messageIds to CHAT selected by URL
  if (url.match(/^\/chat\/\d+\/messages$/)) { 
    console.log("chat/id/messages:", data)
    const id = +url.split("/")[2]
    const chat = client.chats.get(id)
    if (!chat) return;

    const messageIds = data as number[]
    for (const id of messageIds){
      sendWsObject({
        type: "get",
        item: {
          url: `/message/${id}`
        }
      })
    }
    chat.messageIds = messageIds 
  }

  // SET client messages with this format: MessageId => Message 
  if (url.match(/^\/message\/\d+$/)) {
    console.log("message/id", data)
    const id = +url.split("/")[2]
    client.messages.set(id, data as Message)
  }

  // UPDATE list of online users in memory. Dispatch event to render chat list.
  if (url.match(/^\/users\/online$/)) {
    console.log("users/online", data)
    client.usersOnline = data as number[]
    window.dispatchEvent(renderChatList)
  }

}

