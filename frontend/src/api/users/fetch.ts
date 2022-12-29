import axios from "axios"
import { User } from "../../custom"

export const getUserLocal = (id: number) =>
  axios<User | undefined>(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/user/${id}`)
    .then((res) => res.data)
    .catch(() => undefined)
