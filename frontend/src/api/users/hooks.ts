import axios from "axios"
import useSWR from "swr"

import { User } from "@/custom"

const getUser = (id: number) => axios<User>(`/api/user/${id}`).then((res) => res.data)

export const useUser = (id: number | undefined) => {
  const { data, error, mutate } = useSWR<User>(id ? `/api/user/${id}` : null, () =>
    getUser(id as number)
  )

  return { user: data, error, mutate }
}
