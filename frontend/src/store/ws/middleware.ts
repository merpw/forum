import { Middleware } from "redux"

import { RootState } from "@/store/store"
import { wsConnectionActions } from "@/store/ws/connection"
import { chatHandlers } from "@/store/chats"
import { WebSocketResponse } from "@/ws"
import wsActions, { sendGet, sendPost } from "@/store/ws/actions"

const wsConnectionMiddleware: Middleware = (store) => {
  let ws: WebSocket

  return (next) => (action) => {
    if (Object.values(wsConnectionActions).find((a) => a.match(action))) {
      // ignore wsConnection actions
      return next(action)
    }

    if (wsActions.close.match(action)) {
      ws.close()
      store.dispatch(wsConnectionActions.connectionClosed())
      return next(action)
    }

    const token = document.cookie.match(/forum-token=(.*?)(;|$)/)?.[1]
    if (!token) {
      return next(action)
    }

    const state = store.getState() as RootState
    if (state.wsConnection.status === "disconnected") {
      store.dispatch(wsConnectionActions.connectionStarted())

      ws = new WebSocket(`${location.protocol.replace("http", "ws")}//${location.host}/ws`)
      ws.onopen = () => {
        console.log("ws connected")
        ws.send(JSON.stringify({ type: "handshake", item: { token } }))
      }
      ws.onclose = () => {
        console.log("ws disconnected")
        store.dispatch(wsConnectionActions.connectionClosed())
      }

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
    }

    if (sendGet.match(action) || sendPost.match(action)) {
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

    return next(action)
  }
}

export default wsConnectionMiddleware
