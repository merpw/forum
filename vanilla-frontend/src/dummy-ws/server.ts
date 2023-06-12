import WebSocket, { WebSocketServer } from "ws"
import { WebSocketRequest, WSErrorResponse } from "./types"
import getHandler from "./get"
import postHandler from "./post"

const backendUrl =  "http://localhost:8080"

/** Check if session is valid
 * @returns The id of the user
 * @throws Error if session is invalid
 * */
const checkSession = async (token: string) => {
  const resp = await fetch(
    `${backendUrl}/api/internal/check-session?token=${token}`, {
      method: "GET",
      headers: {
        "Internal-Auth": "application/json",
      },
    }
  )

  if (resp.status !== 200) {
    console.error(resp)
    throw new Error("Unexpected response from backend")
  }

  const response = (await resp.json()) as number | { error: string }
  // TODO: maybe change to return { userId: number }

  if (typeof response !== "number") {
    throw new Error(response.error)
  }

  return response
}

const wss = new WebSocketServer({ port: 6969 })

console.log("Server started on ws://localhost:6969")

export const connections = new Map<
  WebSocket.WebSocket,
  {
    userId: number
    token: string
  }
  >()

wss.on("connection", async (ws) => {
  // to show hidden properties
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  console.log("New connection from", (ws as any)._socket.localAddress)
  ws.on("message", async (rawMessage) => {
    try {
      const message: WebSocketRequest<never> = JSON.parse(rawMessage.toString())

      console.log("received:", message)

      if (message.type === "handshake") {
        const userId = await checkSession(message.item.token)

        connections.set(ws, {
          userId,
          token: message.item.token,
        })

        return ws.send(
          JSON.stringify({
            type: "handshake",
            item: {
              data: {
                userId,
              },
            },
          })
        )
      }

      const { userId, token } = connections.get(ws) || {}

      if (!userId || !token) {
        console.log(`Not handshaked, closing connection`)
        return ws.close()
      }

      if ((await checkSession(token)) !== userId) {
        return ws.close()
      }

      const handler = message.type === "get" ? getHandler : postHandler

      const { data } = message.item as never

      return ws.send(
        JSON.stringify(
          handler({
            url: message.item.url,
            data,
            userId,
          })
        )
      )
    } catch (err) {
      if ((err as Error).message === "fetch failed") {
        console.log("Fetch failed, Forum backend is unavailable")
      }
      const errResponse: WSErrorResponse = {
        type: "error",
        item: {
          message: (err as Error).message || "Unknown error",
        },
      }

      ws.send(JSON.stringify(errResponse))
    }
  })
})
