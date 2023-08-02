import { FC } from "react"

import UserLink from "@/components/UserLink"
import Avatar from "@/components/Avatar"
import { useUser } from "@/api/users/hooks"

const ChatInfo: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)

  if (!user) {
    return <div className={"text-info"}>loading...</div>
  }
  return (
    <div className={"flex font-light mb-auto p-2 items-center gap-3"}>
      <Avatar user={user} size={60} className={"w-16"} />
      <div className={"text-info"}>
        Chat with <br />
        <UserLink userId={userId}>
          <span className={"font-Yesteryear gradient-text text-4xl clickable px-1"}>
            {user.username}
          </span>
        </UserLink>
      </div>
    </div>
  )
}

export default ChatInfo
