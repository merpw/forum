import { Group } from "@/api/groups/hooks"

export const getGroupLocal = (id: number): Promise<Group> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/groups/${id}`, {
    headers: {
      "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
    },
  }).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getGroupsLocal = (): Promise<number[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/groups`, {
    headers: {
      "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
    },
  }).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })

export const getGroupMembersLocal = (id: number): Promise<number[]> =>
  fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/groups/${id}/members`, {
    headers: {
      "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
    },
  }).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json()
  })
