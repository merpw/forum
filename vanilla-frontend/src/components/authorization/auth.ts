import { superDivision, client } from "../../main.js"
import { LoginSignup, Index } from "../../pages.js"
import { getMe } from "../../api/get.js"

import { topnavController } from "../feed/topnav.js"
import { categoriesSelector } from "../feed/categories.js"
import { displayPosts } from "../feed/posts.js"

import { ws, wsHandler } from "../chat/ws.js"
import { loginController } from "./login.js"

export const userInfo = {
  Id: -1,
  Name: "",
}

export const Auth = async (session: boolean) => {
  if (session) {
    console.log("client on startup:", client)
    wsHandler()
    superDivision.innerHTML = Index()
    superDivision.classList.replace("login-style", "index-style")
    Object.assign(userInfo, await getMe())
    topnavController()
    categoriesSelector()
    displayPosts("/api/posts")
    window.addEventListener("startChat", client.getAllChats)
  }
  if (!session) {
    if (ws) {
      ws.close()
    }
    console.log("client on close:", client)
    superDivision.innerHTML = LoginSignup()
    superDivision.classList.replace("index-style", "login-style")
    loginController()
  }
}
