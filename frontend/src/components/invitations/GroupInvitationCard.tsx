import { FC } from "react"
import Link from "next/link"

import { Invitation } from "@/api/invitations/hooks"
import { useUser } from "@/api/users/hooks"
import { useGroup } from "@/api/groups/hooks"
import Avatar from "@/components/Avatar"
import { ResponseButtons } from "@/components/invitations/ResponseButtons"

export const GroupInvitationCard: FC<{ invitation: Invitation & { type: 1 } }> = ({
  invitation,
}) => {
  const { user } = useUser(invitation.from_user_id)
  const { group } = useGroup(invitation.associated_id)

  if (!user || !group) return null

  return (
    <div className={"flex flex-col"}>
      <Link href={`/group/${group.id}`} className={"flex gap-2 items-center mr-auto"}>
        <Avatar user={user} size={10} className={"w-12"} />
        <span>
          <span className={"text-primary"}>{user.username}</span> invited you to join their group{" "}
          <span className={"text-primary"}>{group.title}</span>
        </span>
      </Link>
      <ResponseButtons invitation={invitation} />
    </div>
  )
}
