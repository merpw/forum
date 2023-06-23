"use client"

/* eslint-disable import/no-named-as-default-member */
/* https://github.com/iamkun/dayjs/issues/1242 */
import dayjs from "dayjs"
import Link from "next/link"
import { NextPage } from "next/types"
import { FC, useState } from "react"
import relativeTime from "dayjs/plugin/relativeTime"

import { useMe } from "@/api/auth/hooks"
import { useMyPosts } from "@/api/posts/my_posts"
import { useMyPostsLiked } from "@/api/posts/my_posts_liked"
import { PostList } from "@/components/posts/list"
import { User } from "@/custom"
import Avatar from "@/components/Avatar"

/* TODO: add placeholders */

const UserPage: NextPage = () => {
  const { user, isLoading } = useMe()

  const tabs = [
    { title: "Your posts", component: <UserPosts /> },
    { title: "Your liked posts", component: <UserLikedPosts /> },
  ]
  const [activeTab, setActiveTab] = useState(0)

  if (isLoading || !user) return <div>Loading...</div>

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
              <div className={"w-24 sm:w-52"}>
                <Avatar userId={user.id} />
              </div>
              <div className={"self-center text-sm"}>
                <UserInfo user={user} />
              </div>
            </div>
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

/* "2000-01-24" -> "23 years"
 * "2021-01-24" -> "baby 👶"
 */
const calculateAge = (dob: string): string | null => {
  dayjs.extend(relativeTime)

  const parsedDob = dayjs(dob, "YYYY-MM-DD")
  if (!parsedDob.isValid()) return null

  const age = parsedDob.fromNow(true)
  return age.includes("year") ? age + " old" : "baby 👶"
}

const UserInfo: FC<{ user: User }> = ({ user }) => {
  const age = user.dob ? calculateAge(user.dob) : null
  return (
    <>
      <div
        className={"flex flex-col self-center font-light mb-5 text-center sm:text-left text-info"}
      >
        {"Hey, "}
        <span className={"text-3xl sm:text-4xl text-primary font-Yesteryear mx-1"}>
          {user.username}
        </span>
        {"Forgot who you are?"}
      </div>
      <div className={"font-light"}>
        {user.first_name || user.last_name ? (
          <p>
            <span className={"font-light text-info"}>{"Full name"}</span>
            <span
              className={"font-normal start-dot"}
            >{`${user.first_name} ${user.last_name}`}</span>
          </p>
        ) : null}
        {age ? (
          <p>
            <span className={"font-light text-info"}>{"Age"}</span>
            <span className={"font-normal start-dot"}>{age}</span>
          </p>
        ) : null}
        {user.gender ? (
          <p>
            <span className={"font-light text-info"}>{"Gender"}</span>
            <span className={"font-normal start-dot"}>{user.gender}</span>
          </p>
        ) : null}
        <p>
          <span className={"font-light text-info"}>{"Email"}</span>
          <span className={"font-normal start-dot"}>{user.email}</span>
        </p>
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
