import { commentForm } from "../pages.js"

export const displayCommentSection = (id: string) => {
  fetch(`/api/posts/${id}/comments`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      // Add any additional headers if required
    },
  })
    .then((res) => res.json())
    .then((data) => {
      const commentSection = document.getElementById(`CS${id}`)!
      if (commentSection.classList.contains("close")) {
        commentSection.innerHTML = commentForm(id)
        commentSection.classList.replace("close", "open")
        const createPostForm = document.querySelector<HTMLFormElement>(
          `#comment-form-${id}`
        )
        if (createPostForm) {
          new CommentCreator(createPostForm)
        }
      } else {
        commentSection.innerHTML = ""
        commentSection.classList.replace("open", "close")
      }
      return
    })
    .catch((error) => {
      console.error(error)
    })
}

export class CommentCreator {
  private readonly form: HTMLFormElement

  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private onSubmit(event: Event) {
    event.preventDefault()
    const formData: { content: string } = this.getFormData()

    fetch(`/api/posts/${this.form.id.slice(13)}/comment`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    })
      .then((response) => {
        if (response.ok) {
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
