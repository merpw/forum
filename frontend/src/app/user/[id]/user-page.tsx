"use client"

import { NextPage } from "next"
import { FC, useRef } from "react"
import Link from "next/link"
import { SWRConfig } from "swr"

import { User } from "@/custom"
import { PostList } from "@/components/posts/list"
import { followUser, useUser } from "@/api/users/hooks"
import { UserInfo } from "@/components/profiles/UserInfo"
import { useUserPosts } from "@/api/posts/hooks"

const UserPage: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)
  const { posts } = useUserPosts(userId)
  if (!user) return null

  return (
    <>
      <UserInfo user={user} />

      <div className={"flex flex-wrap justify-center gap-3"}>
        <FollowButton userId={user.id} />

        <ChatButton userId={user.id} />
      </div>

      <div className={"mt-5"}>
        <div className={"text-center"}>
          <h2 className={"tab tab-bordered tab-active cursor-default self-center mb-3"}>
            {user.username}
            {"'s posts"}
          </h2>
        </div>
        {posts && <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />}
      </div>
    </>
  )
}

const ChatButton: FC<{ userId: number }> = ({ userId }) => {
  return (
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
  )
}

const FollowButton: FC<{ userId: number }> = ({ userId }) => {
  const { user, mutate } = useUser(userId)

  const popupRef = useRef<HTMLDialogElement>(null)

  const followStatus = user?.follow_status

  if (!user || followStatus === undefined) return null

  const handleFollow = () =>
    followUser(userId).then((value) => {
      mutate({ ...user, follow_status: value })
    })

  return (
    <>
      <button
        className={"button"}
        onClick={() => {
          if (followStatus >= 1) {
            popupRef.current?.showModal()
            return
          }
          handleFollow()
        }}
      >
        {followStatus ? (
          <>
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
                  "M22 10.5h-6m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235v-.11a6.375 6.375 0 0112.75 0v.109A12.318 12.318 0 0110.374 21c-2.331 0-4.512-.645-6.374-1.766z"
                }
              />
            </svg>
            {followStatus === 1 ? "Unfollow" : "Cancel request"}
          </>
        ) : (
          <>
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
                  "M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235v-.11a6.375 6.375 0 0112.75 0v.109A12.318 12.318 0 0110.374 21c-2.331 0-4.512-.645-6.374-1.766z"
                }
              />
            </svg>
            Follow
          </>
        )}
      </button>
      <dialog ref={popupRef} className={"modal"}>
        <form method={"dialog"} className={"modal-box text-center"}>
          <p className={"py-4"}>
            Are you sure that you want to{" "}
            {followStatus === 1 ? `unfollow` : "cancel follow request"}?
          </p>
          <div className={"modal-action justify-center"}>
            {/* if there is a button in form, it will close the modal */}
            <button className={"btn"} onClick={handleFollow}>
              Yes
            </button>
            <button className={"btn btn-ghost"}>Cancel</button>
          </div>
        </form>
      </dialog>
    </>
  )
}

const UserPageWrapper: NextPage<{ user: User }> = ({ user }) => {
  return (
    <SWRConfig
      value={{
        fallback: {
          [`/api/users/${user.id}`]: user,
        },
      }}
    >
      <UserPage userId={user.id} />
    </SWRConfig>
  )
}

export default UserPageWrapper
