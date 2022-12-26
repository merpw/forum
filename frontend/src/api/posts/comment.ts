import { Comment } from "../../custom"
import { user as dummy_user } from "../dummy"

export const CreateComment = (post_id: number, text: string) => {
  return Promise.reject("Not connected to backend yet")
  const comment: Comment = { author: dummy_user, text, date: new Date().toISOString() }
  return Promise.resolve(comment)
}
