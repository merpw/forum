import { postForm } from "../pages.js"
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
  const postFormElement = document.getElementById("create-post")!
  if (postFormElement.classList.contains("close")) {
    // opened.state = true;
    postFormElement.innerHTML = postForm()
    postFormElement.classList.replace("close", "open")
    const createPostForm = document.querySelector<HTMLFormElement>("#post-form")
    if (createPostForm) {
      new PostCreator(createPostForm)
    }
    return
  } else {
    // opened.state = false;
    postFormElement.innerHTML = ""
    postFormElement.classList.replace("open", "close")
    return
  }
}

// Goes to home page
export const goHome = () => {
  const postFormElement = document.getElementById("create-post")!
  // opened.state = false;
  if (postFormElement.classList.contains("open-create-post")) {
    postFormElement.classList.replace("open-create-post", "close-create-post")
  }
  displayPosts("/api/posts")
}

export const topnavController = () => {
  document
    .getElementById("topnav-post")!
    .addEventListener("click", openCloseCreatePost)
  document.getElementById("topnav-home")!.addEventListener("click", goHome)
  document.querySelector(".logout")?.addEventListener("click", logout)
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
        console.log("cookie removed. Logging out.")
        Auth(false)
      }
    })
    .catch((error) => {
      console.error(error)
    })
}
