import { Group } from "@/api/groups/hooks"
import edgeFetcher from "@/api/edge-fetcher"

export const getGroupLocal = (id: number) => edgeFetcher<Group>(`/api/groups/${id}`)

export const getGroupsLocal = () => edgeFetcher<number[]>("/api/groups")

export const getGroupMembersLocal = (id: number) =>
  edgeFetcher<number[]>(`/api/groups/${id}/members`)
