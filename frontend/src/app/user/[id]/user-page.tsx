"use client"

import { NextPage } from "next"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"

import { Post, User } from "@/custom"
import { useMe } from "@/api/auth"
import { PostList } from "@/components/posts/list"

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
      <div>
        <h1 className={"text-2xl mb-5"}>{user.name}</h1>
        <div>
          <h2 className={"text-xl mb-2"}>Recent posts:</h2>
          <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
        </div>
      </div>
    </>
  )
}

export default UserPage
