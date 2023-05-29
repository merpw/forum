import { Auth } from "./components/auth.js"

export const backendUrl = "http://localhost:8080"
export const superDivision = document.getElementById(
  "super-division"
) as HTMLElement

if (document.cookie) {
  Auth(true)
} else {
  Auth(false)
}
