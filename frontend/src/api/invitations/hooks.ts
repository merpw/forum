import useSWR from "swr"

import { dummyFetcher } from "@/api/invitations/dummy"

export type Invitation = {
  id: number
  type: InvitationType
  associatedId: number
  userId: number
  timestamp: string
}

/** 0-following,1-group invitation, 2-group join, 3-event */
export type InvitationType = 0 | 1 | 2 | 3

const fetcher = dummyFetcher

export const useInvitations = () => {
  const { data, error, mutate } = useSWR<number[]>("/api/invitations", fetcher, {
    refreshInterval: 30 * 1000,
  })

  const loading = !data && !error

  return {
    invitations: data,
    loading,
    error,
    mutate,
  }
}

export const useInvitation = (id: number) => {
  const { data, error, mutate } = useSWR<Invitation>(`/api/invitations/${id}`, fetcher)

  const loading = !data && !error

  return {
    invitation: data,
    loading,
    error,
    mutate,
  }
}
