import axios from "axios"
import useSWR from "swr"

import { FollowStatus, User } from "@/custom"
import { useAppSelector } from "@/store/hooks"

const getUser = (id: number) => axios<User>(`/api/users/${id}`).then((res) => res.data)

export const useUser = (id: number | undefined) => {
  const { data, error, mutate } = useSWR<User>(id ? `/api/users/${id}` : null, () =>
    getUser(id as number)
  )

  return { user: data, error, mutate }
}

export const useUsers = () => {
  const { data, error, mutate } = useSWR<number[]>(`/api/users`, () =>
    axios<number[]>(`/api/users`).then((res) => res.data)
  )

  return { users: data, error, mutate }
}

export const useIsUserOnline = (id: number): boolean | undefined => {
  return useAppSelector((state) => state.chats.usersOnline?.includes(id as number))
}

export const useUsersOnline = () => {
  const usersOnline = useAppSelector((state) => state.chats.usersOnline)

  return { usersOnline }
}

export const useUserList = (userIds: number[]) => {
  const { data, error, mutate } = useSWR<User[]>(
    userIds.map((id) => `/api/users/${id}`),
    () => Promise.all(userIds.map((id) => getUser(id)))
  )

  return { users: data, error, mutate }
}

export const togglePrivacy = async (): Promise<boolean | undefined> =>
  axios
    .post<boolean>(`/api/me/privacy`)
    .then((res) => res.data)
    .catch(() => undefined)

export const followUser = async (id: number): Promise<FollowStatus | undefined> =>
  axios
    .post<FollowStatus>(`/api/users/${id}/follow`)
    .then((res) => res.data)
    .catch(() => undefined)
