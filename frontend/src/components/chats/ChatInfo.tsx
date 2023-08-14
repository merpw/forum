import { FC } from "react"
import Link from "next/link"

import UserLink from "@/components/UserLink"
import Avatar from "@/components/Avatar"
import { useUser } from "@/api/users/hooks"
import { useGroup } from "@/api/groups/hooks"
import GroupAvatar from "@/components/groups/Avatar"

const ChatInfo: FC<{ groupId: number } | { userId: number }> = (props) => {
  if ("groupId" in props) {
    return <GroupChatInfo groupId={props.groupId} />
  }
  return <UserChatInfo userId={props.userId} />
}

const UserChatInfo: FC<{ userId: number }> = ({ userId }) => {
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

const GroupChatInfo: FC<{ groupId: number }> = ({ groupId }) => {
  const { group } = useGroup(groupId)

  if (!group) return null

  return (
    <div className={"flex font-light mb-auto p-2 items-center gap-3"}>
      <GroupAvatar group={group} size={60} className={"w-16"} />
      <div className={"text-info"}>
        Group chat
        <br />
        <Link
          href={`/group/${group.id}`}
          className={"font-Yesteryear gradient-text text-4xl clickable px-1"}
        >
          {group.title}
        </Link>
      </div>
    </div>
  )
}

export default ChatInfo
