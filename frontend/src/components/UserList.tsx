import { FC } from "react"
import Link from "next/link"

import { useUser, useUsers, useUsersOnline } from "@/api/users/hooks"
import { useMe } from "@/api/auth/hooks"
import { User } from "@/custom"
import Avatar from "@/components/Avatar"
import { useCollapseIfMobile } from "@/components/chats/section/ChatSection"

const UserList = () => {
  const { users } = useUsers()
  const { usersOnline } = useUsersOnline()

  const { user: reqUser } = useMe()

  if (!users || !usersOnline || !reqUser) {
    return <div className={"text-info"}>loading...</div>
  }

  return (
    <>
      <h1 className={"text-info text-sm my-1"}>Total: {users.length} users</h1>
      <div className={"flex flex-col mt-3 ml-1 gap-3"}>
        <div className={"mt-1"}>
          <MeCard user={reqUser} />
        </div>
        <hr className={"border-dotted border-t-0 border-b-4 border-info opacity-20 my-1"} />
        {users
          .filter((userId) => usersOnline.includes(userId) && userId !== reqUser?.id)
          .map((userId) => (
            <UserCard key={userId} id={userId} />
          ))}

        {users
          .filter((userId) => !usersOnline.includes(userId))
          .map((userId) => (
            <UserCard key={userId} id={userId} />
          ))}
      </div>
    </>
  )
}

const MeCard: FC<{ user: User }> = ({ user }) => {
  const collapseIfMobile = useCollapseIfMobile()

  return (
    <div className={"flex flex-row w-full rounded-3xl hover:bg-neutral hover:saturate-150"}>
      <Link className={"flex clickable w-full self-center"} href={`/me`} onClick={collapseIfMobile}>
        <Avatar user={user} size={30} className={"mr-2 w-7"} />
        <span className={"end-dot text-primary self-center"}>{user.username}</span>
        <span className={"text-info font-light self-center"}>You</span>
      </Link>
    </div>
  )
}

export const UserCard: FC<{ id: number }> = ({ id }) => {
  const { user } = useUser(id)

  const collapseIfMobile = useCollapseIfMobile()

  if (!user) {
    return <div className={"text-info"}>loading...</div>
  }

  return (
    <div
      className={
        "flex flex-row w-full justify-between hover:bg-neutral hover:saturate-150 rounded-3xl"
      }
    >
      <Link
        className={"clickable self-center w-full text-primary flex"}
        href={`/user/${id}`}
        onClick={collapseIfMobile}
      >
        <Avatar user={user} size={30} className={"mr-2 w-7"} />
        <span className={"self-center"}>{user.username}</span>
      </Link>

      <Link
        className={"flex justify-end self-center clickable text-sm text-primary mr-1.5 gap-0.5"}
        href={`/chat/u${id}`}
        onClick={collapseIfMobile}
      >
        <span className={"self-center"}>
          <svg
            xmlns={"http://www.w3.org/2000/svg"}
            fill={"none"}
            viewBox={"0 0 24 24"}
            strokeWidth={1.7}
            stroke={"currentColor"}
            className={"w-5 h-5"}
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
      </Link>
    </div>
  )
}

export default UserList
