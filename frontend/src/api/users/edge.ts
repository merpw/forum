import { User } from "@/custom"

export const getUserLocal = (id: number): Promise<User> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/users/${id}`).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })
