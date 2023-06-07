"use client"

import { NextPage } from "next"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"
import Link from "next/link"

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
      <div className={"hero"}>
        <div className={"hero-content px-0"}>
          <div
            className={
              "card flex-shrink-0 w-full shadow-lg gradient-light dark:gradient-dark px-1 sm:px-3"
            }
          >
            <div className={"card-body sm:flex-row sm:gap-5"}>
              <div tabIndex={0} className={"avatar rounded-full my-1 self-center"}>
                <div className={"w-36 sm:w-48 rounded-full ring-4 ring-neutral"}>
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
              <div className={"self-center sm:text-md font-light text-center sm:text-left"}>
                {/* TODO: add user info if they follows you */}
                {"Hey! I'm "}
                <p className={"text-4xl font-bold text-primary font-Yesteryear mx-1"}>
                  {user.name}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div className={"text-center m-3"}>
        {/* TODO: add valid endpoint */}
        <Link href={"/chat"} className={"button"}>
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
          <span className={"ml-1 text-xs"}>{"Let's chat!"}</span>
        </Link>
      </div>
      <div className={"mt-5"}>
        <div className={"text-center"}>
          <h2 className={"tab tab-bordered tab-active cursor-default self-center mb-3"}>
            {user.name}
            {"'s posts"}
          </h2>
        </div>
        <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
      </div>
    </>
  )
}

export default UserPage
