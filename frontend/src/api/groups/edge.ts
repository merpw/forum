import { Group } from "@/api/groups/hooks"
import { DUMMY_GROUPS } from "@/api/groups/dummy"

export const getGroupLocal = (id: number): Promise<Group> => {
  const group = DUMMY_GROUPS.find((group) => group.id === id)
  if (!group) throw new Error("Group not found")
  return Promise.resolve(group)
}
// fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/groups/${id}`, {
//   headers: {
//     "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
//   },
// }).then((res) => {
//   if (!res.ok) throw new Error("Network response was not ok")
//   return res.json()
// })

export const getGroupsLocal = (): Promise<number[]> =>
  Promise.resolve(DUMMY_GROUPS.map((group) => group.id))

// fetch(`${process.env.FORUM_BACKEND_PRIVATE_URL}/api/groups`, {
//   headers: {
//     "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
//   },
// }).then((res) => {
//   if (!res.ok) throw new Error("Network response was not ok")
//   return res.json()
// })
