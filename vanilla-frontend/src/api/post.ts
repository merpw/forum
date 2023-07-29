import { superDivision } from "../main.js"
import { errorPage } from "../pages.js"
import { CreatePostBody, LoginForm, SignupForm } from "../types"

import { Auth } from "../components/authorization/auth.js"

import { displayPosts, updatePostValues } from "../components/feed/posts.js"
import { displayCommentSection } from "../components/feed/comments.js"
import { openCloseCreatePost } from "../components/feed/topnav.js"

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

const login = async (formData: LoginForm): Promise<boolean | undefined> => {
  return await fetch(`/api/login`, {
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

const signup = async (formData: SignupForm): Promise<boolean | undefined> => {
  return await fetch(`/api/signup`, {
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

const postCreatePost = (formData: CreatePostBody): void => {
  fetch(`/api/posts/create/`, {
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
          alert(error)
        })
      }
    })
    .catch((error) => {
      console.error(error)
    })
}

const postComment = (postId: number, formData: { content: string }): void => {
  fetch(`/api/posts/${postId}/comment`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  })
    .then((r) => {
      if (r.ok) {
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

const likePost = (id: number): void => {
  fetch(`/api/posts/${id}/like`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((r) => {
      if (!r.ok) throw new Error("Something went wrong. Please try again.")
      updatePostValues(id)
    })
    .catch((error) => {
      alert(error)
    })
}

const dislikePost = (id: number): void => {
  fetch(`/api/posts/${id}/dislike`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((r) => {
      if (!r.ok) throw new Error("Something went wrong. Please try again.")
      updatePostValues(id)
    })
    .catch((error) => {
      alert(error)
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
