import { FC } from "react"

import { useIsUserOnline, useUser } from "@/api/users/hooks"
import UserLink from "@/components/UserLink"

const ChatInfo: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)
  const isOnline = useIsUserOnline(userId)

  if (!user) {
    return <div>loading...</div>
  }
  return (
    <h1 className={"text-2xl mb-auto pb-2 border-b"}>
      Chat with{" "}
      <UserLink userId={userId}>
        <span className={"font-bold clickable"}>{user?.username}</span>
        {isOnline ? "🟢" : "🔴"}
      </UserLink>
    </h1>
  )
}

export default ChatInfo
