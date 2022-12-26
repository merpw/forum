import Head from "next/head"
import { useRouter } from "next/router"
import { NextPage } from "next/types"
import { useEffect, useState } from "react"
import { PostList } from "../components/posts"
import { useUser } from "../api/auth"
import { useUserPosts } from "../api/posts/fetch"

const UserPage: NextPage = () => {
  const router = useRouter()
  const { user, isError: isUserError, isLoading: isUserLoading, isLoggedIn } = useUser()
  const { posts, isError: isPostsError, isLoading: isPostsLoading } = useUserPosts(user?.id)

  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects
  useEffect(() => {
    if (!isUserLoading && !isLoggedIn && !isRedirecting) {
      setIsRedirecting(true)
      router.push("/login")
    }
  }, [router, isLoggedIn, isRedirecting, isUserLoading])

  if (isUserError) {
    return <div>Error</div>
  }

  if (isUserLoading) {
    return <div>Loading...</div>
  }

  return (
    <>
      <Head>
        <title>{`Profile - Forum`}</title>
      </Head>
      <h1 className={"text-2xl font-thin mb-5"}>
        {"Hello, "}
        <span className={"text-3xl font-normal"}>{user?.name}</span>
      </h1>
      <h2 className={"text-2xl mb-3"}>Your recent posts:</h2>
      {isPostsLoading && <div>Loading...</div>}
      {isPostsError && <div>Error</div>}
      {posts && <PostList posts={posts} />}
    </>
  )
}

export default UserPage
