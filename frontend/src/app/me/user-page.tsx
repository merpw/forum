"use client"

/* eslint-disable import/no-named-as-default-member */
/* https://github.com/iamkun/dayjs/issues/1242 */
import dayjs from "dayjs"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { NextPage } from "next/types"
import { FC, useEffect, useState } from "react"
import relativeTime from "dayjs/plugin/relativeTime"

import { useMe } from "@/api/auth"
import { useMyPosts } from "@/api/posts/my_posts"
import { useMyPostsLiked } from "@/api/posts/my_posts_liked"
import { PostList } from "@/components/posts/list"
import { User } from "@/custom"

/* TODO: add placeholders */

const UserPage: NextPage = () => {
  const router = useRouter()
  const { user, isLoading, isLoggedIn } = useMe()

  const [isRedirecting, setIsRedirecting] = useState(false) // Prevents duplicated redirects

  const tabs = [
    { title: "Your posts", component: <UserPosts /> },
    { title: "Your liked posts", component: <UserLikedPosts /> },
  ]
  const [activeTab, setActiveTab] = useState(0)

  useEffect(() => {
    if (!isLoading && !isLoggedIn && !isRedirecting) {
      setIsRedirecting(true)
      router.push("/login")
    }
  }, [router, isLoggedIn, isRedirecting, isLoading])

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
              <div tabIndex={0} className={"avatar rounded-full my-1 self-center"}>
                <div className={"w-36 sm:w-48 rounded-full ring-4 ring-accent"}>
                  {/* TODO: add Online/Offline ring. Online: ring-accent; Offline: ring-neutral */}
                  <svg
                    xmlns={"http://www.w3.org/2000/svg"}
                    viewBox={"0 0 24 24"}
                    fill={"currentColor"}
                    className={"opacity-30 w-auto"}
                  >
                    <path
                      d={
                        "M18.685 19.097A9.723 9.723 0 0021.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 003.065 7.097A9.716 9.716 0 0012 21.75a9.716 9.716 0 006.685-2.653zm-12.54-1.285A7.486 7.486 0 0112 15a7.486 7.486 0 015.855 2.812A8.224 8.224 0 0112 20.25a8.224 8.224 0 01-5.855-2.438zM15.75 9a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
                      }
                    />
                  </svg>
                  {/* TODO: add Avatar */}
                </div>
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
      <div className={"flex flex-col self-center text-md font-thin mb-5 text-center sm:text-left"}>
        {"Hey, "}
        <span className={"text-3xl sm:text-4xl font-bold text-primary font-Yesteryear mx-1"}>
          {user.name}
        </span>
        {"Forgot who you are?"}
      </div>
      <div className={"text-md font-thin"}>
        {user.first_name || user.last_name ? (
          <p>
            {"Full name"}
            <span className={"font-light"}> • {`${user.first_name} ${user.last_name}`}</span>
          </p>
        ) : null}
        {age ? (
          <p>
            {"Age"}
            <span className={"font-light"}> • {age}</span>
          </p>
        ) : null}
        {user.gender ? (
          <p>
            {"Gender"}
            <span className={"font-light"}> • {user.gender}</span>
          </p>
        ) : null}
        <p>
          {"Email"}
          <span className={"font-light"}> • {user.email}</span>
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
