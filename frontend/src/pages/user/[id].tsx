import { GetStaticPaths, GetStaticProps, NextPage } from "next"
import { useRouter } from "next/router"
import { FC, useEffect, useRef, useState } from "react"
import { NextSeo } from "next-seo"
import { AxiosError } from "axios"
import Link from "next/link"

import { Post, User } from "@/custom"
import { getUserLocal } from "@/api/users/fetch"
import { useMe } from "@/api/auth"
import { getUserPostsLocal } from "@/api/posts/fetch"
import { PostList } from "@/components/posts/list"
import { useCreateChat, useUserChat } from "@/api/chats/chats"

const UserPage: NextPage<{ user: User; posts: Post[] }> = ({ user, posts }) => {
  const router = useRouter()
  const { user: req_user } = useMe()
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
      <NextSeo title={user.name} description={`Posts created by ${user.name}`} />
      <div>
        <h1 className={"text-2xl mb-3"}>{user.name}</h1>
        <ChatButton userId={user.id} />
        <div>
          <h2 className={"text-xl mb-2"}>Recent posts:</h2>
          <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
        </div>
      </div>
    </>
  )
}

const ChatButton: FC<{ userId: number }> = ({ userId }) => {
  const { chatId } = useUserChat(userId)
  const wasNull = useRef(false)

  const router = useRouter()

  useEffect(() => {
    if (chatId === null) wasNull.current = true
    else if (wasNull.current) router.push(`/chat/${chatId}`)
  }, [chatId, router])

  return (
    <>
      {chatId !== null ? (
        <Link href={`/chat/${chatId}`} className={"clickable text-2xl mb-3"}>
          Open chat
        </Link>
      ) : (
        <CreateChatButton userId={userId} />
      )}
    </>
  )
}

const CreateChatButton: FC<{ userId: number }> = ({ userId }) => {
  const createChat = useCreateChat()
  return (
    <button className={"clickable text-2xl mb-3"} onClick={() => createChat(userId)}>
      Create chat
    </button>
  )
}

export const getStaticPaths: GetStaticPaths<{ id: string }> = async () => {
  return {
    paths: [],
    fallback: "blocking", // fallback tries to regenerate ArtistPage if Artist did not exist during building
  }
}

export const getStaticProps: GetStaticProps<
  { user: User; posts: Post[] },
  { id: string }
> = async ({ params }) => {
  if (!params?.id) {
    return { notFound: true, revalidate: 60 }
  }
  try {
    const user = await getUserLocal(+params.id)
    const posts = await getUserPostsLocal(user.id)

    return {
      props: {
        user,
        posts,
      },
      revalidate: 60,
    }
  } catch (e) {
    if ((e as AxiosError).response?.status !== 404) {
      throw e
    }
    return { notFound: true, revalidate: 60 }
  }
}
export default UserPage
