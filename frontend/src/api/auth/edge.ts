import edgeFetcher from "@/api/edge-fetcher"

const checkSession = async (token: string) => {
  const response = await edgeFetcher<{ error: string } | number>(
    "/api/internal/check-session?token=" + token
  )
  if (typeof response === "object" && response.error) {
    throw new Error(response.error)
  }

  return response as number
}

/** check if the user has permission to see the post */
export const checkPermissions = async (postId: number, userId: number | null) => {
  const response = await edgeFetcher<{ error: string } | boolean>(
    `/api/internal/check-permissions?postId=${postId}` + (userId ? `&userId=${userId}` : "")
  )
  if (typeof response === "object" && response.error) {
    throw new Error(response.error)
  }

  return response as boolean
}

export default checkSession
