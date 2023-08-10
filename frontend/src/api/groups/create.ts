import axios from "axios"

export const createGroup = async (
  title: string,
  description: string,
  invite: number[]
): Promise<number | undefined> =>
  axios
    .post<number>(`/api/groups`, { title, description, invite })
    .then((res) => res.data)
    .catch((err) => {
      throw new Error(err.response?.data?.length < 200 ? err.response.data : "Unexpected error")
    })
