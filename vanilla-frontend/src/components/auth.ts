import { LoginSignup, Index } from "../pages.js"
import { topnavController } from "./topnav.js"
import { superDivision } from "../main.js"
import { categoriesSelector } from "./categories.js"
import { loginController } from "./login.js"
import { displayPosts } from "./posts.js"
import { displayChatUsers } from "./chat.js"

export const Auth = (session: boolean) => {
  if (session) {
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
    // Adding the HTML and changing style
    superDivision.innerHTML = LoginSignup()
    superDivision.classList.replace("index-style", "login-style")
    loginController()
    return
  }
}
