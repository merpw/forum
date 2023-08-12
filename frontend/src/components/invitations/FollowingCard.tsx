import { FC } from "react"
import Link from "next/link"

import { Invitation } from "@/api/invitations/hooks"
import { useUser } from "@/api/users/hooks"
import Avatar from "@/components/Avatar"
import { ResponseButtons } from "@/components/invitations/ResponseButtons"

export const FollowingCard: FC<{ invitation: Invitation & { type: 0 } }> = ({ invitation }) => {
  const { user } = useUser(invitation.from_user_id)

  if (!user) return null

  return (
    <div className={"flex flex-col"}>
      <Link href={`/user/${invitation.from_user_id}`} className={"flex gap-2 items-center mr-auto"}>
        <Avatar user={user} size={10} className={"w-12"} />
        <span>
          <span className={"text-primary"}>{user.username}</span> wants to follow you
        </span>
      </Link>
      <ResponseButtons invitation={invitation} />
    </div>
  )
}
