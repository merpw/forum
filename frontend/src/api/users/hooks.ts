import axios from "axios"
import useSWR from "swr"

import { User } from "@/custom"
import { useAppSelector } from "@/store/hooks"

const getUser = (id: number) => axios<User>(`/api/users/${id}`).then((res) => res.data)

export const useUser = (id: number | undefined) => {
  const { data, error, mutate } = useSWR<User>(id ? `/api/users/${id}` : null, () =>
    getUser(id as number)
  )

  return { user: data, error, mutate }
}

export const useIsUserOnline = (id: number): boolean | undefined => {
  return useAppSelector((state) => state.chats.usersOnline?.includes(id as number))
}
