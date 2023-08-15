import { User } from "@/custom"
import edgeFetcher from "@/api/edge-fetcher"

export const getUserLocal = (id: number) => edgeFetcher<User>(`/api/users/${id}`)
