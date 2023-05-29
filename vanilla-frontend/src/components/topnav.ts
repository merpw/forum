import { superDivision } from "../main.js"
import { errorPage, postForm } from "../pages.js"
import { Auth } from "./auth.js"
import { sendMessage } from "./chat.js"
import { PostCreator, displayPosts } from "./posts.js"

export const postBtn = document.getElementById(
  "topnav-post"
) as HTMLAnchorElement
export const profileBtn = document.getElementById(
  "topnav-profile"
) as HTMLAnchorElement
export const homeBtn = document.getElementById(
  "topnav-home"
) as HTMLAnchorElement

// Opens the create post section of the feed
export const openCloseCreatePost = () => {
  const postFormElement = document.getElementById("create-post") as HTMLElement
  if (!postFormElement) return
  if (postFormElement.classList.contains("close")) {
    postFormElement.innerHTML = postForm()
    postFormElement.classList.replace("close", "open")
    const createPostForm = document.querySelector<HTMLFormElement>("#post-form")
    if (createPostForm) {
      new PostCreator(createPostForm)
    }
    return
  } else {
    postFormElement.innerHTML = ""
    postFormElement.classList.replace("open", "close")
    return
  }
}

// Goes to home page
export const goHome = () => {
  const postFormElement = document.getElementById("create-post") as HTMLElement
  if (postFormElement.classList.contains("open-create-post")) {
    postFormElement.classList.replace("open-create-post", "close-create-post")
  }
  displayPosts("/api/posts")
}

// Logs out the user, deletes the cookie from backend.
const logout = () => {
  fetch("/api/logout", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  }).then((response) => {
    if (response.ok) {
      Auth(false)
    } else {
      superDivision.innerHTML = errorPage(response.status)
      return
    }
  })
}
// Adds event listener to the topNav
export const topnavController = () => {
  const topNavPost = document.getElementById(
      "topnav-post"
    ) as HTMLAnchorElement,
    topNavHome = document.getElementById("topnav-home") as HTMLAnchorElement,
    topNavLogout = document.getElementById(
      "topnav-logout"
    ) as HTMLAnchorElement,
    topNavChat = document.getElementById("topnav-chat") as HTMLAnchorElement,
    chatClose = document.getElementById("chat-close") as HTMLAnchorElement,
    chatList = document.getElementById("chatlist") as HTMLDivElement,
    chatButton = document.getElementById("chat-send") as HTMLElement

  if (
    topNavPost &&
    topNavHome &&
    topNavLogout &&
    topNavChat &&
    chatClose &&
    chatList &&
    chatButton
  ) {
    topNavPost.addEventListener("click", openCloseCreatePost)
    topNavHome.addEventListener("click", goHome)
    topNavLogout.addEventListener("click", logout)
    chatButton.addEventListener("click", sendMessage)
    topNavChat.addEventListener("click", () => {
      chatList.classList.add("chat-list-open")
      chatList.style.width = "200px"
    })
    chatClose.addEventListener("click", () => {
      chatList.classList.remove("chat-list-open")
      chatList.style.width = "0"
    })
  }
}
