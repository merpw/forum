import { FC } from "react"
import Link from "next/link"

import { useIsUserOnline, useUser, useUsers, useUsersOnline } from "@/api/users/hooks"

const UserList = () => {
  const { users } = useUsers()
  const { usersOnline } = useUsersOnline()

  if (!users || !usersOnline) {
    return <div>loading...</div>
  }

  return (
    <div>
      {users
        .slice()
        .sort((a) => (usersOnline.includes(a) ? -1 : 1))
        .map((userId) => (
          <UserCard key={userId} id={userId} />
        ))}
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
