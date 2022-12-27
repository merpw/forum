import { GetStaticPaths, GetStaticProps, NextPage } from "next"
import Head from "next/head"
import { useRouter } from "next/router"
import { useEffect, useState } from "react"
import { Post, User } from "../../custom"

import { getUserLocal, getUsersLocal } from "../../api/users/fetch"
import { useUser } from "../../api/auth"
import { getUserPosts } from "../../api/posts/fetch"
import { PostList } from "../../components/posts/list"

const UserPage: NextPage<{ user: User; posts: Post[] }> = ({ user, posts }) => {
  const router = useRouter()
  const { user: req_user } = useUser()
  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects

  useEffect(() => {
    // Redirect to /me if user is logged in and is viewing their own profile
    if (!isRedirecting && req_user?.id == user.id) {
      setIsRedirecting(true)
      router.replace("/me")
    }
  }, [router, isRedirecting, req_user?.id, user.id])

  return (
    <>
      <Head>
        <title>{`${user.name} - Forum`}</title>
      </Head>
      <div>
        <h1 className={"text-2xl mb-5"}>{user.name}</h1>
        <div>
          <h2 className={"text-xl mb-2"}>Recent posts:</h2>
          <PostList posts={posts} />
        </div>
      </div>
    </>
  )
}

export const getStaticPaths: GetStaticPaths<{ id: string }> = async () => {
  const users = await getUsersLocal()
  return {
    paths: users.map((user) => {
      return { params: { id: user.id.toString() } }
    }),
    // TODO: maybe remove
    fallback: "blocking", // fallback tries to regenerate ArtistPage if Artist did not exist during building
  }
}

export const getStaticProps: GetStaticProps<
  { user: User; posts: Post[] },
  { id: string }
> = async ({ params }) => {
  if (params == undefined) {
    return { notFound: true }
  }
  const user = await getUserLocal(+params.id)
  if (user == undefined) {
    return { notFound: true }
  }
  const { posts } = await getUserPosts(user.id)

  return { props: { user: user, posts: posts } }
}
export default UserPage
