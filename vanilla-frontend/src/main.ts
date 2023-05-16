import { Auth } from "./components/auth.js"
export const superDivision = document.getElementById(
  "super-division"
) as HTMLElement

if (document.cookie) {
  Auth(true)
} else {
  Auth(false)
}
