import axios from "axios"

import { GroupMemberStatus } from "@/api/groups/hooks"

export const inviteToGroup = async (groupId: string, userId: number) =>
  axios.post(`/api/groups/${groupId}/invite`, { user_id: userId }).then((res) => res.data)

export const joinGroup = async (groupId: string) =>
  axios.post<GroupMemberStatus>(`/api/groups/${groupId}/join`).then((res) => res.data)

export const leaveGroup = async (groupId: string) =>
  axios.post<GroupMemberStatus>(`/api/groups/${groupId}/leave`).then((res) => res.data)
