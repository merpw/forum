import { LoginSignup, Index } from "../pages.js"
import { topnavController } from "./topnav.js"
import { superDivision } from "../main.js"
import { categoriesSelector } from "./categories.js"
import { loginController } from "./login.js"
import { displayPosts } from "./posts.js"
import { ws, wsHandler } from "./ws.js"

import { getMe } from "../api/get.js"

export const userInfo = {
  Id: -1,
  Name: "",
}

export const Auth = async (session: boolean) => {
  if (session) {
    wsHandler()
    setTimeout(async () => {
      // Adding the HTML and changing style
      superDivision.innerHTML = Index()
      superDivision.classList.replace("login-style", "index-style")
      Object.assign(userInfo, await getMe()) // Sets the userInfo (Id, Name)
      topnavController() // Adds event listeners to the top-navigation bar
      categoriesSelector()
      displayPosts("/api/posts")
    }, 100)
  }
  if (!session) {
    if (ws) {
      ws.close()
    }
    // Adding the HTML and changing style
    setTimeout(() => {
      superDivision.innerHTML = LoginSignup()
      superDivision.classList.replace("index-style", "login-style")
      loginController()
    }, 100)
  }
}
