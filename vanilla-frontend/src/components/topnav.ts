import { superDivision } from "../main.js"
import { errorPage, postForm } from "../pages.js"
import { Auth } from "./auth.js"
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
  if (!postFormElement) return;
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

export const topnavController = () => {
  const topNavPost = document.getElementById("topnav-post") as HTMLElement
  const topNavHome = document.getElementById("topnav-post") as HTMLElement
  const topNavLogout = document.querySelector(".logout") as HTMLElement
  if (topNavPost && topNavHome && topNavLogout){
    topNavPost.addEventListener("click", openCloseCreatePost)
    topNavHome.addEventListener("click", goHome)
    topNavLogout.addEventListener("click", logout)
  }
}

// Logs out the user, deletes the cookie from backend.
const logout = () => {
  fetch("/api/logout", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((response) => {
      if (response.ok) {
        // WS[0].close(1000, "User logging out. Closing connection.")
        // WS.pop()
        Auth(false)
      } else {
        superDivision.innerHTML = errorPage(response.status)
        return
      }
    })
}
