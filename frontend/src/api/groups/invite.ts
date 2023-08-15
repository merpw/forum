import axios from "axios"

import { GroupMemberStatus } from "@/api/groups/hooks"

export const inviteToGroup = async (groupId: number, userId: number) =>
  axios
    .post(`/api/groups/${groupId}/invite`, { user_id: userId })
    .then((res) => res.data)
    .catch((err) => {
      throw new Error(err.response.data)
    })

export const joinGroup = async (groupId: number) =>
  axios
    .post<GroupMemberStatus>(`/api/groups/${groupId}/join`)
    .then((res) => res.data)
    .catch((err) => {
      throw new Error(err.response.data)
    })

export const leaveGroup = async (groupId: number) =>
  axios
    .post<GroupMemberStatus>(`/api/groups/${groupId}/leave`)
    .then((res) => res.data)
    .catch((err) => {
      throw new Error(err.response.data)
    })
