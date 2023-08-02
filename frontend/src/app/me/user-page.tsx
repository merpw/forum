"use client"

import Link from "next/link"
import { NextPage } from "next/types"
import { useState } from "react"

import { useMe } from "@/api/auth/hooks"
import { useMyPosts } from "@/api/posts/my_posts"
import { useMyPostsLiked } from "@/api/posts/my_posts_liked"
import { PostList } from "@/components/posts/list"
import Avatar from "@/components/Avatar"
import { UserInfo } from "@/components/profiles/UserInfo"

/* TODO: add placeholders */

const UserPage: NextPage = () => {
  const { user, isLoading } = useMe()

  const tabs = [
    { title: "Your posts", component: <UserPosts /> },
    { title: "Your liked posts", component: <UserLikedPosts /> },
  ]
  const [activeTab, setActiveTab] = useState(0)

  if (isLoading || !user) return <div className={"text-info text-center mt-5 mb-7"}>Loading...</div>

  return (
    <>
      <div className={"hero"}>
        <div className={"hero-content px-0"}>
          <div
            className={
              "card flex-shrink-0 w-full shadow-lg gradient-light dark:gradient-dark px-1 sm:px-3"
            }
          >
            <div className={"card-body sm:flex-row sm:gap-5"}>
              <Avatar user={user} size={50} className={"w-24 sm:w-52 self-center p-1"} />
              <div className={"self-center text-sm"}>
                <UserInfo user={user} />
              </div>
            </div>

            {user.bio && (
              <div className={"mb-5 text-center"}>
                <div className={"font-light text-info start-dot end-dot mb-1"}>About me</div>
                <div className={"text-sm"}>{user.bio}</div>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className={"text-center m-3"}>
        <Link href={"/create"} className={"button"}>
          <span className={"my-auto"}>
            <svg
              xmlns={"http://www.w3.org/2000/svg"}
              fill={"none"}
              viewBox={"0 0 24 24"}
              strokeWidth={1.5}
              stroke={"currentColor"}
              className={"w-5 h-5"}
            >
              <path
                strokeLinecap={"round"}
                strokeLinejoin={"round"}
                d={
                  "M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
                }
              />
            </svg>
          </span>
          <span className={"ml-1 text-xs"}>Create a new post</span>
        </Link>
      </div>
      <div className={"mt-5 space-y-3"}>
        <div className={"text-center"}>
          <ul className={"tab tab-lg p-0"}>
            {tabs.map(({ title }, key) => (
              <li key={key}>
                <button
                  type={"button"}
                  key={key}
                  className={"tab tab-bordered space-y-5 " + (activeTab == key ? "tab-active" : "")}
                  onClick={() => setActiveTab(key)}
                >
                  {title}
                </button>
              </li>
            ))}
          </ul>
        </div>
        <div>{tabs[activeTab].component}</div>
      </div>
    </>
  )
}

const UserPosts = () => {
  const { posts } = useMyPosts()

  if (posts == undefined) return null

  if (posts.length == 0)
    return (
      <div className={"text-info text-center mt-5 mb-7"}>{"You haven't posted anything yet"}</div>
    )

  return <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
}

const UserLikedPosts = () => {
  const { posts } = useMyPostsLiked()

  if (posts == undefined) return null

  if (posts.length == 0)
    return (
      <div className={"text-info text-center mt-5 mb-7"}>{"You haven't liked any posts yet"}</div>
    )

  return <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
}

export default UserPage
