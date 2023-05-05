import axios from "axios"
import useSWR from "swr"

import { User } from "@/custom"

const getUser = (id: number) => axios<User>(`/api/users/${id}`).then((res) => res.data)

export const useUser = (id: number | undefined) => {
  const { data, error, mutate } = useSWR<User>(id ? `/api/users/${id}` : null, () =>
    getUser(id as number)
  )

  return { user: data, error, mutate }
}
