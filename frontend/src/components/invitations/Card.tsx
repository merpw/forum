import { FC } from "react"

import { useInvitation } from "@/api/invitations/hooks"
import { GroupInvitationCard } from "@/components/invitations/GroupInvitationCard"
import { FollowingCard } from "@/components/invitations/FollowingCard"
import { GroupJoinCard } from "@/components/invitations/GroupJoinCard"
import { EventCard } from "@/components/invitations/EventCard"

// TODO: add more types

const Card: FC<{ id: number }> = ({ id }) => {
  const { invitation } = useInvitation(id)

  if (!invitation) return null

  switch (invitation.type) {
    case 0:
      return <FollowingCard invitation={invitation} />
    case 1:
      return <GroupInvitationCard invitation={invitation} />
    case 2:
      return <GroupJoinCard invitation={invitation} />
    case 3:
      return <EventCard invitation={invitation} />
    default:
      return null
  }
}

export default Card
