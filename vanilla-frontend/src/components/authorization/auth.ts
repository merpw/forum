import { superDivision } from "../../main.js"
import { LoginSignup, Index } from "../../pages.js"
import { getMe } from "../../api/get.js"

import { topnavController } from "../feed/topnav.js"
import { categoriesSelector } from "../feed/categories.js"
import { displayPosts } from "../feed/posts.js"

import { ws, wsHandler } from "../chat/ws.js"
import { displayChatUsers } from "../chat/chat.js"

import { loginController } from "./login.js"

export const userInfo = {
  Id: -1,
  Name: "",
}

export const Auth = (session: boolean): void => {
  if (session) {
    wsHandler()
    setTimeout(async () => {
      superDivision.innerHTML = Index()
      superDivision.classList.replace("login-style", "index-style")
      Object.assign(userInfo, await getMe())
      topnavController()
      categoriesSelector()
      displayPosts("/api/posts")
      displayChatUsers()
    }, 100)
  }
  if (!session) {
    if (ws) {
      ws.close()
    }
    setTimeout(() => {
      superDivision.innerHTML = LoginSignup()
      superDivision.classList.replace("index-style", "login-style")
      loginController()
    }, 100)
  }
}
