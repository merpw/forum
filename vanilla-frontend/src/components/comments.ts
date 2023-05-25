import { commentForm } from "../pages.js"
import { updatePostValues } from "./posts.js"

export const displayCommentSection = (id: string) => {
  // Opens and closes comment section if you press the comment section button
  fetch(`/api/posts/${id}/comments`)
    .then((resp) => resp.json())
    .then((comments) => {
      if (!comments) return
      const commentSection = document.getElementById(`CS${id}`) as HTMLElement
      if (!commentSection) {
        return // Add error handling here
      }
      // If statement for opening the comment section
      if (commentSection.classList.contains("close")) {
        // Appends the form to the commentSection
        const commentFormElement = document.createElement("div")
        commentFormElement.className = "comment-form"
        commentFormElement.innerHTML = commentForm(id)
        commentSection.appendChild(commentFormElement)

        commentSection.classList.replace("close", "open")
        const createPostForm = document.querySelector<HTMLFormElement>(
          `#comment-form-${id}`
        )
        if (createPostForm) {
          new CommentCreator(createPostForm)
        }
        // This loops through all the comments and creates them in the DOM.
        comments.reverse()
        for (let i = 0; i < comments.length; i++) {
          // Parent element
          const comment = document.createElement("div")
          comment.className = "comment"
          comment.id = `CommentID${comments[i].id}`
          // Comment-info Child
          const commentInfo = document.createElement("div")
          commentInfo.className = "comment-info"
          const date = new Date(comments[i].date)
          const formatDate = date.toLocaleString("en-GB", { timeZone: "EET" })
          commentInfo.textContent = `${comments[i].author.name}\n\tat ${formatDate}`

          // Comment content Child
          const commentContent = document.createElement("div")
          commentContent.className = "comment-content"
          commentContent.textContent = `${comments[i].content}`

          comment.append(commentInfo, commentContent)
          commentSection.appendChild(comment)
        }
        // Else statement for closing the comment section
      } else {
        commentSection.replaceChildren()
        commentSection.classList.replace("open", "close")
      }
      return
    })
    .catch((error) => {
      console.error(error)
    })
}

//   <div id="CommentID${commentID}" class="comment">
//      <div class="comment-info">
//        <h4>Author</h4>
//      </div>
//      <div class="comment-content"></div>
//   </div>
// `
export class CommentCreator {
  private readonly form: HTMLFormElement

  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private onSubmit(event: Event) {
    event.preventDefault()
    const formData: { content: string } = this.getFormData()
    const postID = this.form.id.slice(13)

    fetch(`/api/posts/${postID}/comment`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    })
      .then((response) => {
        if (response.ok) {
          console.log("PostID in CommentCreator", postID)
          const commentSection = document.getElementById(
            `CS${postID}`
          ) as HTMLElement
          commentSection.replaceChildren()
          commentSection.classList.replace("open", "close")
          updatePostValues(postID)
          displayCommentSection(postID)
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

  private getFormData(): { content: string } {
    const content =
      this.form.querySelector<HTMLInputElement>("#comment-content")
    if (content) {
      return { content: content.value }
    }
    throw new Error("Could not find form input fields.")
  }
}
