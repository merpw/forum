import { Middleware } from "redux"

import { RootState } from "@/store/store"
import { wsConnectionActions } from "@/store/ws/connection"
import { chatHandlers } from "@/store/chats"
import { WebSocketResponse } from "@/ws"
import wsActions, { sendGet, sendPost } from "@/store/ws/actions"

const wsConnectionMiddleware: Middleware = (store) => {
  let ws: WebSocket
  let retryTimeout = 0

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

    if (wsActions.connect.match(action)) {
      store.dispatch(wsConnectionActions.connectionStarted())

      ws = new WebSocket(`${location.protocol.replace("http", "ws")}//${location.host}/ws`)
      ws.onopen = () => {
        console.log("ws connected")
        retryTimeout = 0
        ws.send(JSON.stringify({ type: "handshake", item: { token } }))
      }
      ws.onclose = () => {
        console.log("ws disconnected")
        setTimeout(() => {
          retryTimeout += 1000
          store.dispatch(wsConnectionActions.connectionClosed())
          store.dispatch(wsActions.connect())
        }, retryTimeout)
      }

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as WebSocketResponse<never>

          if (data.type === "handshake") {
            // TODO: maybe use userId from the handshake response
            console.log("ws handshake success")
            return next(wsConnectionActions.connectionEstablished())
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

      return next(action)
    }

    if (state.wsConnection.status === "disconnected") {
      store.dispatch(wsActions.connect())
      return next(action)
    }

    if (sendGet.match(action) || sendPost.match(action)) {
      const type = action.payload.type
      const url = action.payload.item.url

      if (type === "get" && state.wsConnection.pendingRequests.includes(url)) {
        // ignore duplicate requests
        return next(action)
      }

      if (state.wsConnection.status === "connected") {
        ws.send(JSON.stringify(action.payload))
        store.dispatch(wsConnectionActions.requestPending(action.payload.item.url))
      } else {
        setTimeout(() => {
          store.dispatch({ type: "ws/send", payload: action.payload })
        }, retryTimeout)
        return next(action)
      }
    }

    return next(action)
  }
}

export default wsConnectionMiddleware
