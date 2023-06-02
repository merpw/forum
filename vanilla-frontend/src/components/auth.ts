import { LoginSignup, Index } from "../pages.js"
import { topnavController } from "./topnav.js"
import { superDivision } from "../main.js"
import { categoriesSelector } from "./categories.js"
import { loginController } from "./login.js"
import { displayPosts } from "./posts.js"
import { displayChatUsers } from "./chat.js"
import { ws, wsHandler } from "./ws.js"

export const Auth = async (session: boolean) => {
  if (session) {
    if (!ws) {
      wsHandler()
    }
    // Adding the HTML and changing style
    superDivision.innerHTML = Index()
    superDivision.classList.replace("login-style", "index-style")

    topnavController() // Adds event listeners to the top-navigation bar
    categoriesSelector() // TODO: Actual functionality for this
    displayChatUsers()
    displayPosts("/api/posts")
    return
  }
  if (!session) {
    if (ws) {
      ws.close(1000, "User logging out. Closing connection.")
    }
    // Adding the HTML and changing style
    superDivision.innerHTML = LoginSignup()
    superDivision.classList.replace("index-style", "login-style")
    loginController()
    return
  }
}
