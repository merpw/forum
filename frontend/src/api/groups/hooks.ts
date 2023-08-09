import useSWR from "swr"

import { DUMMY_GROUPS } from "@/api/groups/dummy"

export type Group = {
  id: number
  title: string
  description: string

  member_status: GroupMemberStatus
  member_count: number

  creator_id: number
}

/* 0-not member, 1-member, 2-requested */
export type GroupMemberStatus = 0 | 1 | 2

export const useGroups = () => {
  const { data, error, mutate } = useSWR<number[]>("/api/groups", () =>
    DUMMY_GROUPS.map((group) => group.id)
  )

  return {
    groups: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}

export const useGroup = (id: number) => {
  const { data, error, mutate } = useSWR<Group | undefined>(`/api/groups/${id}`, () =>
    DUMMY_GROUPS.find((group) => group.id === id)
  )

  return {
    group: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}
