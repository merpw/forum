import { NextPage } from "next"
import { FC } from "react"
import Link from "next/link"

import { Post, User } from "@/custom"
import { PostList } from "@/components/posts/list"
import Avatar from "@/components/Avatar"

const UserPage: NextPage<{ user: User; posts: Post[] }> = ({ user, posts }) => {
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
              <Avatar userId={user.id} className={"w-24 sm:w-48 m-auto self-center"} />
              <div className={"self-center font-light text-center sm:text-left"}>
                {/* TODO: add user info if they follows you */}
                {"Hey! I'm "}
                <p className={"text-4xl text-primary font-Yesteryear mx-1"}>{user.username}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <ChatButton userId={user.id} />

      <div className={"mt-5"}>
        <div className={"text-center"}>
          <h2 className={"tab tab-bordered tab-active cursor-default self-center mb-3"}>
            {user.username}
            {"'s posts"}
          </h2>
        </div>
        <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
      </div>
    </>
  )
}

const ChatButton: FC<{ userId: number }> = ({ userId }) => {
  return (
    <div className={"text-center m-3"}>
      <Link href={`/chat/u${userId}`}>
        <button className={"button"}>
          <span className={"my-auto"}>
            <svg
              xmlns={"http://www.w3.org/2000/svg"}
              fill={"none"}
              viewBox={"0 0 24 24"}
              strokeWidth={1.5}
              stroke={"currentColor"}
              className={"w-5 h-5 scale-x-[-1]"}
            >
              <path
                strokeLinecap={"round"}
                strokeLinejoin={"round"}
                d={
                  "M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 01-2.555-.337A5.972 5.972 0 015.41 20.97a5.969 5.969 0 01-.474-.065 4.48 4.48 0 00.978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25z"
                }
              />
            </svg>
          </span>
          <span className={"text-xs"}>{"Let's chat!"}</span>
        </button>
      </Link>
    </div>
  )
}

export default UserPage
