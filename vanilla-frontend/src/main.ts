import { Auth } from "./components/auth.js"
export const superDivision = document.getElementById(
  "super-division"
) as HTMLElement

/* Not optimal solution, very temporary fix. */
export const WS: WebSocket[] = []

/* 
 * COMMENTED OUT WS LOGIC:
 * 39-40 in login.ts (adds websocket)
 * 68-69 in topnav.ts (closes websocket)
 * NEED TO IMPORT WS IN BOTH.
 */

if (document.cookie) {
  Auth(true)
} else {
  Auth(false)
}
