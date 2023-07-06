import { postForm } from "../pages.js"
import { PostCreator, displayPosts } from "./posts.js"
import { logout } from "../api/post.js"
import { displayChatUsers } from "./chat.js"

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
export const openCloseCreatePost = async () => {
  const postFormElement = document.getElementById("create-post") as HTMLElement
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
  const categories = document.querySelectorAll(".category-title") as NodeListOf <HTMLElement>
  if (postFormElement.classList.contains("open-create-post")) {
    postFormElement.classList.replace("open-create-post", "close-create-post")
  }
  categories.forEach((category) => {
    if (category.classList.contains("selected")) {
      category.classList.remove("selected")
      return
    }
  })
  displayPosts("/api/posts")
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
    chatDisplay = document.getElementById("chatlist") as HTMLDivElement,
    chat = document.getElementById("chat-area") as HTMLDivElement

  topNavPost.addEventListener("click", openCloseCreatePost)
  topNavHome.addEventListener("click", goHome)
  topNavLogout.addEventListener("click", logout)
  topNavChat.addEventListener("click", () => {
    setTimeout(() => {
      chatDisplay.classList.add("chat-list-open")
      chatDisplay.style.width = "200px"
      chat.style.left = "210px"
    }, 50)
  })
  chatClose.addEventListener("click", () => {
    chatDisplay.classList.remove("chat-list-open")
    chatDisplay.style.width = "0"
    chat.style.left = "10px"
  })
}
