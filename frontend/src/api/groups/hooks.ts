import useSWR from "swr"

import fetcher from "@/api/fetcher"

export type Group = {
  id: number
  title: string
  description: string

  member_status?: GroupMemberStatus
  member_count: number

  creator_id: number
}

/* 0-not member, 1-member, 2-requested */
export type GroupMemberStatus = 0 | 1 | 2

export const useGroups = () => {
  const { data, error, mutate } = useSWR<number[]>("/api/groups", fetcher)

  return {
    groups: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}

export const useGroup = (id: number) => {
  const { data, error, mutate } = useSWR<Group | undefined>(`/api/groups/${id}`, fetcher)

  return {
    group: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}

export const useGroupMembers = (id: number, withPedning?: boolean) => {
  const { data, error, mutate } = useSWR<number[]>(
    `/api/groups/${id}/members${withPedning ? "?withPending" : ""}`,
    fetcher
  )

  return {
    members: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}
