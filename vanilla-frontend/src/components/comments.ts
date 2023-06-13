import { commentForm } from "../pages.js"
import { getComments } from "../api/get.js"
import { postComment } from "../api/post.js"
import { createElement, iterator } from "./utils.js"

export const displayCommentSection = async (id: string) => {
  const comments: iterator = await getComments(id)
  if (!comments) return

  // Opens and closes comment section if you press the comment section button
  const commentSection = document.getElementById(`CS${id}`) as HTMLElement

  // If statement for opening the comment section
  if (commentSection.classList.contains("close")) {
    // Appends the form to the commentSection
    const commentFormElement = createElement(
      "div",
      "comment-form",
      null,
      null,
      commentForm(id)
    ) as HTMLDivElement
    commentSection.appendChild(commentFormElement)
    commentSection.classList.replace("close", "open")

    const createPostForm = document.querySelector(
      `#comment-form-${id}`
    ) as HTMLFormElement
    new CommentCreator(createPostForm)

    // Loop through all the comments and creates them in the DOM.
    for (const c of Object.values(comments)) {
      const date = new Date(c.date)
      const formatDate = date.toLocaleString("en-GB", { timeZone: "EET" })
      // Parent element
      const comment = createElement(
        "div",
        "comment",
        `CommentID${c.id}`
      ) as HTMLDivElement
      // Comment info
      const commentInfo = createElement(
        "div",
        "comment-info",
        null,
        `${c.author.name}\n\tat ${formatDate}`
      )
      // Comment content
      const commentContent = createElement(
        "div",
        "comment-content",
        null,
        `${c.content}`
      )
      comment.append(commentInfo, commentContent)
      commentSection.appendChild(comment)
    }
    // Else statement for closing the comment section
  } else {
    commentSection.replaceChildren()
    commentSection.classList.replace("open", "close")
  }
  return
}

export class CommentCreator {
  private readonly form: HTMLFormElement
  constructor(form: HTMLFormElement) {
    this.form = form
    this.form.addEventListener("submit", this.onSubmit.bind(this))
  }

  private onSubmit(event: Event) {
    event.preventDefault()
    const postId = this.form.id.slice(13)
    const formData: { content: string } = this.getFormData()
    postComment(postId, formData)
  }

  private getFormData(): { content: string } {
    const content = this.form.querySelector(
      "#comment-content"
    ) as HTMLInputElement
    return { content: content.value }
  }
}
