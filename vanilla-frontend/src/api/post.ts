import { Auth } from "../components/auth.js"
import { superDivision, backendUrl } from "../main.js"
import { errorPage } from "../pages.js"
import { loginController } from "../components/login.js"
import { CreatePostBody, LoginForm, SignupForm } from "../types"
import { displayPosts, updatePostValues } from "../components/posts.js"
import { displayCommentSection } from "../components/comments.js"
import { openCloseCreatePost } from "../components/topnav.js"

// Logs out the user, deletes the cookie from backend.
const logout = async () => {
  await fetch("/api/logout", {
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

const login = async (formData: LoginForm) => {
  await fetch(`${backendUrl}/api/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  }).then((response) => {
    if (response.ok) {
      Auth(true)
    } else {
      response.text().then((error) => {
        console.log(`Error: ${error}`)
      })
    }
  })
}

const signup = async (formData: SignupForm) => {
  await fetch(`${backendUrl}/api/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  })
  .then((response) => {
    if (response.ok) {
      loginController()
    } else {
      response.text().then((error) => {
        console.log(`Error: ${error}`)
      })
    }
  })
  .catch((error) => {
    console.error(error)
  })

}

const postCreatePost = (formData: CreatePostBody) => {
fetch(`${backendUrl}/api/posts/create/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    })
      .then((response) => {
        if (response.ok) {
          openCloseCreatePost()
          displayPosts("/api/posts/")
          return
        } else {
          response.text().then((error) => {
            console.log(`Error: ${error}`)
            // TODO: Displaying error message to user.
          })
        }
      })
      .catch((error) => {
        console.error(error)
      })

}

const postComment = (postId: string, formData: unknown) => {
  fetch(`${backendUrl}/api/posts/${postId}/comment`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  })
    .then((response) => {
      if (response.ok) {
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
        response.text().then((error) => {
          console.log(`Error: ${error}`)
        })
      }
    })
    .catch((error) => {
      console.error(error)
    })
}

export { login, signup, logout, postComment, postCreatePost }
