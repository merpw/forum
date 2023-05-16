import { CreatePostBody, SafePost } from "POST"
import { openCloseCreatePost } from "./topnav.js"
import { displayCommentSection } from "./comments.js"

export const displayPosts = (endpoint: string) => {
  const postsDisplay = document.getElementById("posts-display") as HTMLElement
  postsDisplay.innerHTML = ""
  const postList: SafePost[] = [];
  // const postBuffer: SafePost[] = [];
  fetch(endpoint)
    .then((response) => response.json())
    .then((posts) => {
      // Posts is the array of objects sent by the server
      for (const post of posts) {
        const PostData: SafePost = {
          HTML:  formattedPost (
            post.id.toString(),
            post.title,
            post.author.name,
            post.author.id.toString(),
            post.categories,
            post.comments_count.toString(),
            post.likes_count.toString(),
            post.dislikes_count.toString()
          ),
          Content: post.content
        }
        postList.push(PostData)
      }

      postList.reverse()

      // TODO: some kind of buffer for the posts. Can and will be reused for chat.
      // for (const post of postList) {
      // }

      for (const post of postList) {
        const postElement = document.createElement("div")
        postElement.className = "post-wrapper"
        postElement.innerHTML = post.HTML;

        postsDisplay.appendChild(postElement)
      }

      // Event listener for the like button
      const likeButton = document.querySelectorAll(
        ".post-likes"
      ) as NodeListOf<HTMLElement>
      likeButton.forEach((button) => {
        button.addEventListener("click", () => {
          fetch(`/api/posts/${button.id.slice(1)}/like`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
          })
            .then((res) => {
              if (!res.ok) throw new Error()
              // Call the function to update the post values after liking
              updatePostValues(button.id.slice(1))
            })
            .catch((error) => {
              console.error(error)
              // Handle the error if the request fails
            })
        })
      })
      const dislikeButton = document.querySelectorAll(
        ".post-dislikes"
      ) as NodeListOf<HTMLElement>
      dislikeButton.forEach((button) => {
        button.addEventListener("click", () => {
          fetch(`/api/posts/${button.id.slice(1)}/dislike`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
          })
            .then((res) => {
              if (!res.ok) throw new Error()
              // Call the function to update the post values after liking
              updatePostValues(button.id.slice(1))
            })
            .catch((error) => {
              console.error(error)
              // Handle the error if the request fails
            })
        })
      })
      const commentButton = document.querySelectorAll(
        ".post-comments"
      ) as NodeListOf<HTMLElement>
      commentButton.forEach((button) => {
        button.addEventListener("click", () => {
          displayCommentSection(button.id.slice(1))
          updatePostValues(button.id.slice(1))
        })
      })
    })
}

const updatePostValues = (postId: string) => {
  fetch(`/api/posts/${postId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((response) => response.json())
    .then((data) => {
      // Extract the updated values from the response data
      const updatedCommentsCount = data.comments_count
      const updatedLikesCount = data.likes_count
      const updatedDislikesCount = data.dislikes_count

      // Update the post UI elements with the new values
      const postCommentsElement = document.getElementById(`C${postId}`)
      const postLikesElement = document.getElementById(`L${postId}`)
      const postDislikesElement = document.getElementById(`D${postId}`)

      if (postCommentsElement && postLikesElement && postDislikesElement) {
        postCommentsElement.innerHTML = `<i class='bx bx-comment'></i> ${updatedCommentsCount}`
        postLikesElement.innerHTML = `<i class='bx bx-like'></i> ${updatedLikesCount}`
        postDislikesElement.innerHTML = `<i class='bx bx-dislike'></i> ${updatedDislikesCount}`
      }
    })
    .catch((error) => {
      console.error(error)
      // Handle the error if the request fails
    })
}

export class PostCreator {
  private readonly form: HTMLFormElement

  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private onSubmit(event: Event) {
    event.preventDefault()
    const formData: CreatePostBody = this.getFormData()

    fetch("/api/posts/create/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    })
      .then((response) => {
        if (response.ok) {
          openCloseCreatePost()
          displayPosts("/api/posts")
          return
          // TODO: Something after post is created. Maybe close post window?
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

  private getFormData(): CreatePostBody {
    // Form inputs
    const title = this.form.querySelector<HTMLInputElement>("#post-title"),
      content = this.form.querySelector<HTMLInputElement>("#post-content"),
      category = this.form.querySelector<HTMLInputElement>("#post-category")

    if (title && content && category) {
      const formData: CreatePostBody = {
        Title: title.value,
        Content: content.value,
        Description: `${title.value} ${category.value} ${content.value}`,
        Categories: [category.value],
      }
      return formData
    }
    throw new Error("Could not find form input fields.")
  }
}

const formattedPost = (
  id: string,
  title: string,
  author: string,
  authorId: string,
  category: string,
  commentCount: string,
  likeCount: string,
  dislikeCount: string
): string => {
  return `
  <div class="post" id="${id}">
	<div class="post-information">
		<div class="post-title"><h4>${title}</h4></div>
		<div class="post-author" id="${authorId}">by ${author}</div>
		<div class="post-categories">#${category}</div>
		<hr>
	</div>

	<div class="post-content">
		
	</div>

	<div class="post-footer" id="${id}">
		<div class="post-comments post-icon" id="C${id}"><i class='bx bx-comment'></i> ${commentCount}</div>
		<div class="post-likes post-icon" id="L${id}"><i class='bx bx-like' ></i> ${likeCount}</div>
		<div class="post-dislikes post-icon" id="D${id}"><i class='bx bx-dislike' ></i> ${dislikeCount}</div>
	</div>
  <section class="comments-section close" id=CS${id}></section>
  </div>
	`
}
