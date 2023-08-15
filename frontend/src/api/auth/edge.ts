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
export default checkSession
