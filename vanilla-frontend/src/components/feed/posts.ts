import { CreatePostBody } from "../../types"

import { getPosts, getPostValues } from "../../api/get.js"
import { dislikePost, likePost, postCreatePost } from "../../api/post.js"

import { displayCommentSection } from "./comments.js"

/* Utilities */
import { createElement, iterator } from "../utils.js"

export const displayPosts = async (endpoint: string): Promise<void> => {
  const postsDisplay = document.getElementById("posts-display") as HTMLElement
  postsDisplay.replaceChildren()
  const postList: Array<HTMLDivElement> = []
  const posts: iterator = await getPosts(endpoint)

  // Posts is the array of objects sent by the server
  for (const post of Object.values(posts)) {
    const date = new Date(post.date)
    const formatDate = date.toLocaleString("en-GB", { timeZone: "EET" })

    const postDiv = formattedPost(
      post.id,
      post.title,
      post.author.username,
      post.author.id,
      formatDate,
      post.categories,
      post.description,
      post.comments_count,
      post.likes_count,
      post.dislikes_count
    )

    postList.unshift(postDiv)
  }

  for (const post of postList) {
    postsDisplay.appendChild(post)
  }
}

// Updates likes/dislikes/comments when values change
export const updatePostValues = async (postId: number): Promise<void> => {
  const post: iterator = await getPostValues(postId),
    postCommentsElement = document.getElementById(`C${postId}`) as HTMLElement,
    postLikesElement = document.getElementById(`L${postId}`) as HTMLElement,
    postDislikesElement = document.getElementById(`D${postId}`) as HTMLElement

  postCommentsElement.innerHTML = `<i class='bx bx-comment'></i> ${post.comments_count}`
  postLikesElement.innerHTML = `<i class='bx bx-like'></i> ${post.likes_count}`
  postDislikesElement.innerHTML = `<i class='bx bx-dislike'></i> ${post.dislikes_count}`
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
    postCreatePost(formData)
  }

  // Gets all values from the form, and puts it in CreatePostBody type.
  private getFormData(): CreatePostBody {
    const title = this.form.querySelector("#post-title") as HTMLInputElement,
      content = this.form.querySelector("#post-content") as HTMLInputElement,
      category = this.form.querySelector("#post-category") as HTMLInputElement

    return {
      Title: title.value,
      Content: content.value,
      Description: content.value,
      Categories: [category.value],
    }
  }
}

// Creates a post in the DOM and returns the HTMLDivElement
const formattedPost = (
  id: number,
  title: string,
  author: string,
  authorId: number,
  date: string,
  category: string,
  content: string,
  commentCount: number,
  likeCount: number,
  dislikeCount: number
): HTMLDivElement => {
  const newPost = createElement("div", "post", `PostId${id}`) as HTMLDivElement

  // Creates the post information div
  const postInformation = createElement("div", "post-information")
  postInformation.addEventListener("click", (e) => {
    displayCommentSection(id)
    updatePostValues(id)
    e.preventDefault()
  })

  // Creates the post title
  const postTitle = createElement("div", "post-title")
  const titleContent = createElement("h4", null, null, title)
  postTitle.appendChild(titleContent)
  // Creates the author and time div
  const postAuthor = createElement(
    "div",
    "post-author",
    `AuthorId${authorId}`,
    `by ${author} at ${date}`
  )

  // Creates the category div
  const postCategories = createElement(
    "div",
    "post-categories",
    null,
    `#${category}`
  )

  postInformation.append(
    postTitle,
    postAuthor,
    postCategories,
    createElement("hr")
  )

  const postContent = createElement("div", "post-content", null, content)
  postContent.addEventListener("click", (e) => {
    displayCommentSection(id)
    updatePostValues(id)
    e.preventDefault()
  })

  const postFooter = createElement("div", "post-footer")
  const postComments = createElement(
    "div",
    "post-comments post-icon",
    `C${id}`,
    null,
    `<i class="bx bx-comment" style="font-size: 20px; margin-right: 5px;"></i>  ${commentCount}`
  )
  postComments.addEventListener("click", () => {
    displayCommentSection(id)
    updatePostValues(id)
  })

  const postLikes = createElement(
    "div",
    "post-likes post-icon",
    `L${id}`,
    null,
    `<i class="bx bx-like" style="font-size: 20px; margin-right: 5px;"></i>  ${likeCount}`
  )

  const postDislikes = createElement(
    "div",
    "post-dislikes post-icon",
    `D${id}`,
    null,
    `<i class="bx bx-dislike" style="font-size: 20px; margin-right: 5px;"></i>  ${dislikeCount}`
  )

  postLikes.addEventListener("click", () => likePost(id))
  postDislikes.addEventListener("click", () => dislikePost(id))
  postFooter.append(postComments, postLikes, postDislikes)

  const commentSection = createElement(
    "section",
    "comments-section close",
    `CS${id}`
  )

  newPost.append(postInformation, postContent, postFooter, commentSection)

  return newPost
}
