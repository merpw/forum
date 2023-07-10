import { FC } from "react"
import Link from "next/link"

import { useIsUserOnline, useUser, useUsers, useUsersOnline } from "@/api/users/hooks"
import { useMe } from "@/api/auth/hooks"
import { User } from "@/custom"

const UserList = () => {
  const { users } = useUsers()
  const { usersOnline } = useUsersOnline()

  const { user: reqUser } = useMe()

  if (!users || !usersOnline || !reqUser) {
    return <div>loading...</div>
  }

  return (
    <div>
      <MeCard user={reqUser} />
      {users
        .filter((userId) => usersOnline.includes(userId) && userId !== reqUser?.id)
        .map((userId) => (
          <UserCard key={userId} id={userId} />
        ))}
      <hr />
      {users
        .filter((userId) => !usersOnline.includes(userId))
        .map((userId) => (
          <UserCard key={userId} id={userId} />
        ))}
    </div>
  )
}

const MeCard: FC<{ user: User }> = ({ user }) => {
  return (
    <div>
      <Link className={"clickable"} href={`/me`}>
        {user.username}
      </Link>
      🟢
    </div>
  )
}

const UserCard: FC<{ id: number }> = ({ id }) => {
  const { user } = useUser(id)
  const isUserOnline = useIsUserOnline(id)

  if (!user) {
    return <div>loading...</div>
  }

  return (
    <div>
      <Link className={"clickable"} href={`/user/${id}`}>
        {user.username}
      </Link>
      {isUserOnline ? "🟢" : "🔴"}
      <Link href={`/chat/u${id}`}>chat</Link>
    </div>
  )
}

export default UserList
