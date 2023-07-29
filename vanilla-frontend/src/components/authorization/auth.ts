import { superDivision, client } from "../../main.js"
import { LoginSignup, Index } from "../../pages.js"
import { getMe, getUserById, getUserIds } from "../../api/get.js"

import { topnavController } from "../feed/topnav.js"
import { categoriesSelector } from "../feed/categories.js"
import { displayPosts } from "../feed/posts.js"

import { ws, wsHandler } from "../chat/ws.js"
import { loginController } from "./login.js"
import { User } from "../../types.js"

interface forumState {
  me: User
  users: Map <number, User>
 }


export const state: forumState = {
  me: {
    Id: -1,
    Name: ""
  },
  users: new Map<number, User>(),
}

async function initializeState() {
    Object.assign(state.me, await getMe())
    const ids = await getUserIds()
    const filteredIds = ids.filter(id => id !== state.me.Id)
    for(const id of filteredIds){
      state.users.set(id, await getUserById(id))
    }
}

export const Auth = async (session: boolean) => {
  if (session) {
    await initializeState()

    superDivision.innerHTML = Index()
    superDivision.classList.replace("login-style", "index-style")

    const greeting = document.getElementById("greeting-name") as HTMLElement
    greeting.textContent = state.me.Name

    topnavController()
    categoriesSelector()
    displayPosts("/api/posts")

    wsHandler()
  }

  if (!session) {
    if (ws) {
      ws.close()
    }
    client.reset()
    superDivision.innerHTML = LoginSignup()
    superDivision.classList.replace("index-style", "login-style")
    loginController()
  }
}
