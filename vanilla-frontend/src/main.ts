import { Auth } from "./components/authorization/auth.js"
import { Client } from "./components/chat/client.js"

// export const backendUrl = "http://localhost:8080"
export const superDivision = document.getElementById(
  "super-division"
) as HTMLElement

export let client = new Client()


if (document.cookie) {
  Auth(true)
} else {
  Auth(false)
}
