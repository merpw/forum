import { Auth } from "./components/auth.js"
export const superDivision = document.getElementById(
  "super-division"
) as HTMLElement // superDivision, basically the html body tag, but better.

if (document.cookie) {
  Auth(true)
} else {
  Auth(false)
}
