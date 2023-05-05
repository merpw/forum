import { Middleware } from "redux"
import { createAction } from "@reduxjs/toolkit"

import { RootState } from "@/store/store"
import { wsConnectionActions } from "@/store/ws"
import { chatHandlers } from "@/store/chats/chats"
import { WebSocketResponse } from "@/ws"

// TODO: use FORUM_BACKEND_URL
const WebSocketUrl = "ws://localhost:6969"

const wsConnectionMiddleware: Middleware = (store) => {
  let ws: WebSocket

  return (next) => (action) => {
    const state = store.getState() as RootState

    if (
      state.wsConnection.status === "disconnected" &&
      action.type !== wsConnectionActions.connectionStarted.type
    ) {
      store.dispatch(wsConnectionActions.connectionStarted())
      ws = new WebSocket(WebSocketUrl)
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as WebSocketResponse<never>

          if (data.type === "handshake") {
            // TODO: maybe use userId from the handshake response
            store.dispatch(wsConnectionActions.connectionEstablished())
            console.log("ws handshake success")
            return
          }

          if (data.type === "error") {
            console.error("ws error", data.item.message)
            return
          }

          if (!data.item?.url) {
            console.error("invalid ws response", data)
            return
          }

          const { handler } = chatHandlers.find(({ regex }) => data.item.url.match(regex)) ?? {}

          if (!handler) {
            console.error("unhandled ws response", data)
            return
          }

          store.dispatch(wsConnectionActions.requestDone(data.item.url))
          store.dispatch(handler(data as never))
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
        store.dispatch(wsConnectionActions.connectionClosed())
      }
    }

    if (action.type === "ws/send") {
      const type = action.payload.type
      const url = action.payload.item.url

      if (type === "get" && state.wsConnection.pendingRequests.includes(url)) {
        return action
      }

      if (state.wsConnection.status === "connected") {
        ws.send(JSON.stringify(action.payload))
        store.dispatch(wsConnectionActions.requestPending(action.payload.item.url))
      } else {
        setTimeout(() => {
          store.dispatch({ type: "ws/send", payload: action.payload })
        }, 100)
        return action
      }
    }

    if (action.type === "ws/close") {
      ws.close()
      store.dispatch(wsConnectionActions.connectionClosed())
    }

    return next(action)
  }
}

export const sendWSGet = createAction("ws/send", (url: string) => {
  return {
    payload: {
      type: "get",
      item: {
        url,
      },
    },
  }
})

// TODO: improve types, maybe even infer data type by url
export const sendWSPost = createAction("ws/send", <T>(url: string, data: T) => {
  return {
    payload: {
      type: "post",
      item: {
        url,
        data,
      },
    },
  }
})

export default wsConnectionMiddleware
