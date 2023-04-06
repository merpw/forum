import axios from "axios"
import { User } from "../../custom"

export const getUserLocal = (id: number) =>
  axios<User>(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/user/${id}`).then((res) => res.data)
