import { FC } from "react"

import { useIsUserOnline, useUser } from "@/api/users/hooks"
import UserLink from "@/components/UserLink"
import Avatar from "@/components/Avatar"

const ChatInfo: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)
  const isOnline = useIsUserOnline(userId)

  if (!user) {
    return <div className={"text-info text-center mt-5 mb-7"}>loading...</div>
  }
  return (
    <div className={"flex font-light mb-auto p-2 items-center gap-3"}>
      <div className={"pt-2 w-16"}>
        <Avatar userId={userId} />
      </div>
      <div className={"text-info"}>
        Chat with <br />
        <UserLink userId={userId}>
          <span className={"font-Yesteryear gradient-text text-4xl clickable"}>
            {user?.username}
          </span>
        </UserLink>
      </div>
    </div>
  )
}

export default ChatInfo
