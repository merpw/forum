import { Post } from "../../custom"
import { posts, user as dummy_user } from "../dummy"

export const CreatePost = (title: string, content: string) => {
  return Promise.reject("Not connected to backend yet")
  const post: Post = {
    title,
    content,
    id: posts.length + 1,
    likes: 0,
    dislikes: 0,
    author: dummy_user,
    date: new Date().toISOString(),
    comments: [],
  }
  posts.push(post)
  return Promise.resolve(post.id)
}
