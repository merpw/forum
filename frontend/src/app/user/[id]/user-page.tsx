import { NextPage } from "next"
import { FC } from "react"
import Link from "next/link"

import { Post, User } from "@/custom"
import { PostList } from "@/components/posts/list"

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
