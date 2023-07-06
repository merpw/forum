const checkSession = async (token: string) => {
  const response = await fetch(
    `${process.env.FORUM_BACKEND_PRIVATE_URL}/api/internal/check-session?token=${token}`,
    {
      headers: {
        "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
      },
    }
  )

  const body = await response.json()
  if (body.error) {
    throw new Error(body.error)
  }

  return body as number
}
export default checkSession
