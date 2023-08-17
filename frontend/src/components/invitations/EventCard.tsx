import { FC } from "react"
import Link from "next/link"

import { Invitation } from "@/api/invitations/hooks"
import { ResponseButtons } from "@/components/invitations/ResponseButtons"
import { useEvent } from "@/api/events/hooks"
import { useUser } from "@/api/users/hooks"
import { useGroup } from "@/api/groups/hooks"
import Avatar from "@/components/Avatar"

export const EventCard: FC<{ invitation: Invitation & { type: 3 } }> = ({ invitation }) => {
  const { event } = useEvent(invitation.associated_id)

  const { user } = useUser(invitation.from_user_id)

  if (!event || !user) return null

  return (
    <div className={"flex flex-col"}>
      <Link href={`/group/${event.group_id}/events`} className={"flex gap-2 items-center mr-auto"}>
        <Avatar user={user} size={10} className={"w-12"} />
        <span>
          <span className={"text-primary"}>{user.username}</span> created an event{" "}
          <span className={"text-primary"}>{event.title}</span> in{" "}
          <span className={"text-primary"}>
            <GroupTitle id={event.group_id} />
          </span>
        </span>
      </Link>
      <ResponseButtons invitation={invitation} acceptText={"I'm in!"} declineText={"I'm out"} />
    </div>
  )
}

const GroupTitle: FC<{ id: number }> = ({ id }) => {
  const { group } = useGroup(id)

  if (!group) return null

  return <>{group.title}</>
}
