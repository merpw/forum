import { superDivision } from "../../main.js"
import { LoginSignup, Index } from "../../pages.js"
import { getMe } from "../../api/get.js"

import { topnavController } from "../feed/topnav.js"
import { categoriesSelector } from "../feed/categories.js"
import { displayPosts } from "../feed/posts.js"

import { ws, wsHandler } from "../chat/ws.js"
import { displayChatUsers } from "../chat/chat.js"

import { loginController } from "./login.js"

// userInfo keeps info of the authorized user
export const userInfo = {
  Id: -1,
  Name: "",
}

// Auth keeps state of wether or not the user is authorized or not
export const Auth = (session: boolean): void => {
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
      displayChatUsers()
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
