import { WSGetResponse, WSPostResponse, WebSocketResponse, Message } from "../../types"

import { Auth } from "../authorization/auth.js"

import { iterator } from "../utils.js"

export const messageEvent = new Event("messageEvent")

export let ws: WebSocket

export const wsHandler = (): void => {
  let retryTimeout = 0
  ws = new WebSocket(`${location.protocol.replace("http", "ws")}//${location.host}/ws`)
  ws.onmessage = (event: MessageEvent): void => {
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
    retryTimeout = 0
    ws.send(JSON.stringify({ type: "handshake", item: { token } }))
  }

  ws.onclose = () => {
    console.log("ws disconnected")
    setTimeout(() => {
      retryTimeout += 1000
      wsHandler() 
    }, retryTimeout)
  }}

const postHandler = async (resp: WSPostResponse<never>): Promise<void> => {
  const data = resp.item.data
  const url = resp.item.url

  if (url.match(/^\/chat\/\d+\/message$/)) { // /chat/{id}/message

  }

  if (url.match(/^\/chat\/create$/)) { // /chat/create
    
  }
}

const getHandler = (resp: WSGetResponse<never>) => {
  const data: iterator = resp.item.data
  const url = resp.item.url

  /* /chat/{id} */
  if (url.match(/^\/chat\/\d+$/)) {

  }

  /* /chat/{id}/messages */
  if (url.match(/^\/chat\/\d+\/messages$/)) { 

  }

  /* /message/{id} */
  if (url.match(/^\/message\/\d+$/)) {

  }
 /* /users/online */
  if (url.match(/^\/users\/online$/)) {

  }

  /* /chat/all */
  if (url.match(/^\/chat\/all$/)) {

  }
}

