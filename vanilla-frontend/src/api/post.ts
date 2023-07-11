/* IMPORTS */

/* Root */
import { superDivision, backendUrl } from "../main.js"
import { errorPage } from "../pages.js"
import { CreatePostBody, LoginForm, SignupForm } from "../types"

/* Authorization */
import { Auth } from "../components/authorization/auth.js"

/* Feed */
import { displayPosts, updatePostValues } from "../components/feed/posts.js"
import { displayCommentSection } from "../components/feed/comments.js"
import { openCloseCreatePost } from "../components/feed/topnav.js"

// Logs out the user, deletes the cookie from backend.
const logout = async (): Promise<void> => {
  await fetch("/api/logout", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  }).then((r) => {
    if (r.ok) {
      Auth(false)
    } else {
      superDivision.innerHTML = errorPage(r.status)
      return
    }
  })
}

// Logs in the user. Adds cookie from backend.
const login = async (formData: LoginForm): Promise<boolean | undefined> => {
  return await fetch(`${backendUrl}/api/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  }).then((r) => {
    if (r.ok) {
      return true
    } else {
      r.text().then((error) => {
        alert(error)
        return false
      })
    }
  })
}

// Signs up the user, calls login() if successful signup
const signup = async (formData: SignupForm): Promise<boolean | undefined> => {
  return await fetch(`${backendUrl}/api/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  })
  .then((r) => {
    if (r.ok) {
      return true
    } else {
      r.text().then((error) => {
        alert(error)
          return false
      })
    }
  })
  .catch((error) => {
    alert(error)
      return false
  })
}

// Creates a post and sends a post request to backend with the form data
const postCreatePost = (formData: CreatePostBody): void => {
  fetch(`${backendUrl}/api/posts/create/`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  })
    .then((r) => {
      if (r.ok) {
        openCloseCreatePost()
        displayPosts("/api/posts/")
        return
      } else {
        r.text().then((error) => {
          alert(error)})        
      }
    })
    .catch((error) => {
      console.error(error)
    })
}

// Comments on a post and sends a post request to the backend with the formData
const postComment = (postId: string, formData: unknown): void => {
  fetch(`${backendUrl}/api/posts/${postId}/comment`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  })
    .then((r) => {
      if (r.ok) {
        console.log("PostID in CommentCreator", postId)
        const commentSection = document.getElementById(
          `CS${postId}`
        ) as HTMLElement
        commentSection.replaceChildren()
        commentSection.classList.replace("open", "close")
        updatePostValues(postId)
        displayCommentSection(postId)
        return
      } else {
        r.text().then((error) => {
          alert(error)
        })
      }
    })
    .catch((error) => {
      console.error(error)
    })
}

// Likes a post and sends a request to the backend
const likePost = (id: string): void => {
  fetch(`${backendUrl}/api/posts/${id}/like`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((r) => {
      if (!r.ok) throw new Error("Something went wrong. Please try again.")
      // Call the function to update the post values after liking
      updatePostValues(id)
    })
    .catch((error) => {
        alert(error)
      // Handle the error if the request fails
    })
}

// Dislikes a post and sends a request to the backend
const dislikePost = (id: string): void => {
  fetch(`${backendUrl}/api/posts/${id}/dislike`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((r) => {
      if (!r.ok) throw new Error("Something went wrong. Please try again.")
      // Call the function to update the post values after liking
      updatePostValues(id)
    })
    .catch((error) => {
        alert(error)
      // Handle the error if the request fails
    })
}

export {
  login,
  signup,
  logout,
  postComment,
  postCreatePost,
  likePost,
  dislikePost,
}
