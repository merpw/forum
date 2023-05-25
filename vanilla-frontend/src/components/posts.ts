import { CreatePostBody } from "./types"
import { openCloseCreatePost } from "./topnav.js"
import { displayCommentSection } from "./comments.js"

export const displayPosts = (endpoint: string) => {
  const postsDisplay = document.getElementById("posts-display") as HTMLElement
  postsDisplay.replaceChildren()
  const postList: Array<HTMLDivElement> = [];
  fetch(endpoint)
    .then((response) => response.json())
    .then((posts) => {
      // Posts is the array of objects sent by the server
      for (const post of posts) {
          const date = new Date(post.date)
          const formatDate = date.toLocaleString('en-GB', { timeZone: 'EET' })
  
       const postDiv = formattedPost (
            post.id.toString(),
            post.title,
            post.author.name,
            post.author.id.toString(),
            formatDate,
            post.categories,
            post.content,
            post.comments_count.toString(),
            post.likes_count.toString(),
            post.dislikes_count.toString()
          )
        postList.push(postDiv)
      }
      postList.reverse()
      // TODO: some kind of buffer for the posts. Can and will be reused for chat.
      // for (const post of postList) {
      // }
      for (const post of postList) {
        postsDisplay.appendChild(post)
      }
    })
}

export const updatePostValues = (postId: string) => {
  fetch(`/api/posts/${postId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((response) => response.json())
    .then((data) => {
      // Update the post UI elements with the new values
      const postCommentsElement = document.getElementById(`C${postId}`)
      const postLikesElement = document.getElementById(`L${postId}`)
      const postDislikesElement = document.getElementById(`D${postId}`)

      if (postCommentsElement && postLikesElement && postDislikesElement) {
        postCommentsElement.innerHTML = `<i class='bx bx-comment'></i> ${data.comments_count}`
        postLikesElement.innerHTML = `<i class='bx bx-like'></i> ${data.likes_count}`
        postDislikesElement.innerHTML = `<i class='bx bx-dislike'></i> ${data.dislikes_count}`
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
          displayPosts("/api/posts/")
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
  
  // Gets all values from the form, and puts it in CreatePostBody type.
  // CreatePostBody {
  //  Title:       string
  //  Content:     string
  //  Description: string
  //  Categories:  string[]
  // }
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
  dateAndTime: string,
  category: string,
  content: string,
  commentCount: string,
  likeCount: string,
  dislikeCount: string
): HTMLDivElement => {
  const newPost = document.createElement("div")
  newPost.className = "post"
  newPost.id = `${id}`

  // Creates the post information div
  const postInformation = document.createElement("div")
  postInformation.classList.add("post-information")

  // Creates the post title
  const postTitle = document.createElement("div")
  postTitle.className = "post-title"
  const titleContent = document.createElement("h4")
  titleContent.textContent = `${title}`
  postTitle.appendChild(titleContent)

  // Creates the author and time div
  const postAuthor = document.createElement("div")
  postAuthor.className = "post-author"
  postAuthor.id = `${authorId}`
  postAuthor.textContent = `by ${author} at ${dateAndTime}`
  
  // Creates the category div 
  const postCategories = document.createElement("div")
  postCategories.className = "post-categories"
  postCategories.textContent = `#${category}` 

  postInformation.append(postTitle, postAuthor, postCategories, document.createElement("hr"))
  
  // Creates post-content 
  const postContent = document.createElement("div")
  postContent.className = "post-content"
  postContent.textContent = `${content}`

  // Creates postFooter
  const postFooter = document.createElement("div")
  postFooter.className = "post-footer"
  postFooter.id = `${id}`

  // Creates comments count div
  const postComments = document.createElement("div")
  postComments.className = "post-comments post-icon"
  postComments.id = `C${id}`
  postComments.innerHTML = `<i class="bx bx-comment"></i> ${commentCount}`
  postComments.addEventListener("click", () => {
    displayCommentSection(id)
    updatePostValues(id)
  })

  const postLikes = document.createElement("div")
  postLikes.className = "post-likes post-icon"
  postLikes.id = `L${id}`
  postLikes.innerHTML = `<i class="bx bx-like"></i> ${likeCount}`
  postLikes.addEventListener("click", () => {
    fetch(`/api/posts/${id}/like`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error()
        // Call the function to update the post values after liking
        updatePostValues(id)
      })
      .catch((error) => {
        console.error(error)
        // Handle the error if the request fails
      })
  })

  const postDislikes = document.createElement("div")
  postDislikes.className = "post-dislikes post-icon"
  postDislikes.id = `D${id}`
  postDislikes.innerHTML = `<i class="bx bx-dislike"></i> ${dislikeCount}`
  postDislikes.addEventListener("click", () => {
    fetch(`/api/posts/${id}/dislike`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error()
        // Call the function to update the post values after liking
        updatePostValues(id)
      })
      .catch((error) => {
        console.error(error)
        // Handle the error if the request fails
      })
  })
  postFooter.append(postComments, postLikes, postDislikes)

  const commentSection = document.createElement("section")
  commentSection.className = "comments-section close"
  commentSection.id = `CS${id}`

  newPost.append(postInformation, postContent, postFooter, commentSection)

  return newPost  
 //  return `
 //  <div class="post" id="${id}">
 //    <div class="post-information">
 //      <div class="post-title"><h4>${title}</h4></div>
 //      <div class="post-author" id="${authorId}">by ${author} at ${date}</div>
 //      <div class="post-categories">#${category}</div>
 //      <hr>
 //  	</div>

	// <div class="post-content">

	// </div>

	// <div class="post-footer" id="${id}">
	// 	<div class="post-comments post-icon" id="C${id}"><i class='bx bx-comment'></i> ${commentCount}</div>
	// 	<div class="post-likes post-icon" id="L${id}"><i class='bx bx-like' ></i> ${likeCount}</div>
	// 	<div class="post-dislikes post-icon" id="D${id}"><i class='bx bx-dislike' ></i> ${dislikeCount}</div>
	// </div>
 //  <section class="comments-section close" id=CS${id}></section>
 //  </div>
	// `
}
