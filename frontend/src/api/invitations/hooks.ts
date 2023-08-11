import useSWR from "swr"
import axios from "axios"

export type Invitation = {
  id: number
  type: InvitationType
  from_user_id: number
  to_user_id: number
  timestamp: string
} & (
  | {
      type: 0
    }
  | {
      type: 1
      /** group id */
      associated_id: number
    }
  | {
      type: 2
      /** group id */
      associated_id: number
    }
  | {
      type: 3
      /** event id */
      associated_id: number
    }
)

/** 0-following,1-group invitation, 2-group join, 3-event */
export type InvitationType = 0 | 1 | 2 | 3

const fetcher = (url: string) => axios.get(url).then((res) => res.data)

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
