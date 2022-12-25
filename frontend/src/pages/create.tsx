import { NextPage } from "next"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import { useUser } from "../fetch/user"

const CreatePost: NextPage = () => {
  const { isLoggedIn } = useUser()
  const router = useRouter()

  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects
  useEffect(() => {
    if (isLoggedIn && !isRedirecting) {
      setIsRedirecting(true)
      router.push("/login")
    }
  }, [router, isLoggedIn, isRedirecting])

  return <div>Create post</div>
}

export default CreatePost
