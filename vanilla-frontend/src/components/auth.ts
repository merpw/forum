import { LoginSignup, Index } from "../pages.js"
import { topnavController } from "./topnav.js"
import { superDivision } from "../main.js"
import { categoriesSelector } from "./categories.js"
import { loginController } from "./login.js"
import { displayPosts } from "./posts.js"
import { displayChatUsers } from "./chat.js"
import { ws, wsHandler } from "./ws.js"
import { getMe } from "../api/get.js"

export const userInfo = {
  Id: -1,
  Name: ""
}

export const Auth = async (session: boolean) => {
  if (session) {
    if (!ws) {
      wsHandler()
    }
    // Adding the HTML and changing style
    superDivision.innerHTML = Index()
    superDivision.classList.replace("login-style", "index-style")

    Object.assign(userInfo, await getMe()) // Sets the userInfo (Id, Name)
    topnavController() // Adds event listeners to the top-navigation bar
    categoriesSelector() 
    displayChatUsers()
    displayPosts("/api/posts")
    return
  }
  if (!session) {
    if (ws){
      ws.close()
    }
    // Adding the HTML and changing style
    superDivision.innerHTML = LoginSignup()
    superDivision.classList.replace("index-style", "login-style")
    loginController()
    return
  }
}
