import { useRouter } from "next/router"
import { NextPage } from "next/types"
import { useEffect, useState } from "react"
import { useUser } from "../fetch/user"

const UserPage: NextPage = () => {
  const router = useRouter()
  const { user, isError, isLoading, isLoggedIn } = useUser()

  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects
  useEffect(() => {
    if (!isLoggedIn && !isRedirecting) {
      setIsRedirecting(true)
      router.push("/login")
    }
  }, [router, isLoggedIn, isRedirecting])

  if (isError) {
    return <div>Error</div>
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <div>
      <h1 className={"text-xl"}>{user?.name}</h1>
    </div>
  )
}

export default UserPage
