import useSWR from "swr"

import fetcher from "@/api/fetcher"

export type Event = {
  id: number
  group_id: number
  created_by: number
  title: string
  description: string
  time_and_date: string
  timestamp: string

  status?: EventStatus
}

// 0 - not going, 1 - going, 2 - requested
export type EventStatus = 0 | 1 | 2

export const useGroupEvents = (groupId: number) => {
  const { data, error, mutate } = useSWR<number[]>(`/api/groups/${groupId}/events`, fetcher)

  return {
    events: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}

export const useEvent = (id: number) => {
  const { data, error, mutate } = useSWR<Event | undefined>(`/api/events/${id}`, fetcher)

  return {
    event: data,
    isLoading: !error && !data,
    error,
    mutate,
  }
}
