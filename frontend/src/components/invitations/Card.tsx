import { FC } from "react"
import Link from "next/link"
import { useSWRConfig } from "swr"

import { Invitation, useInvitation } from "@/api/invitations/hooks"
import { useUser } from "@/api/users/hooks"
import Avatar from "@/components/Avatar"
import respondToInvitation from "@/api/invitations/respond"

// TODO: add more types

const Card: FC<{ id: number }> = ({ id }) => {
  const { invitation } = useInvitation(id)

  if (!invitation) return null

  switch (invitation.type) {
    case 0:
      return <FollowingCard invitation={invitation} />
    default:
      return <div>not implemented yet, type: {invitation.type}</div>
  }
}

const FollowingCard: FC<{ invitation: Invitation }> = ({ invitation }) => {
  const { user } = useUser(invitation.user_id)

  if (!user) return null

  return (
    <div className={"flex flex-col"}>
      <Link href={`/user/${invitation.user_id}`} className={"flex gap-2 items-center"}>
        <span className={"w-12"}>
          <Avatar user={user} size={10} className={""} />
        </span>
        <span>
          <span className={"text-primary"}>{user.username}</span> wants to follow you
        </span>
      </Link>
      <ResponseButtons invitationId={invitation.id} />
    </div>
  )
}

const ResponseButtons: FC<{ invitationId: number }> = ({ invitationId }) => {
  const { mutate } = useSWRConfig()

  const respond = (accept: boolean) =>
    respondToInvitation(invitationId, accept)
      .catch(console.error)
      .finally(() => {
        mutate(`/api/invitations/${invitationId}`)
        mutate(`/api/invitations`)
      })

  return (
    <div className={"flex gap-2"}>
      <button className={"btn btn-sm btn-success"} onClick={() => respond(true)}>
        Accept
      </button>
      <button className={"btn btn-sm btn-error"} onClick={() => respond(false)}>
        Decline
      </button>
    </div>
  )
}

export default Card
