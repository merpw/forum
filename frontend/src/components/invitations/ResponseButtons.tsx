import { FC } from "react"
import { useSWRConfig } from "swr"

import respondToInvitation from "@/api/invitations/respond"
import { Invitation } from "@/api/invitations/hooks"

export const ResponseButtons: FC<{
  invitation: Invitation
  acceptText?: string
  declineText?: string
}> = ({ invitation, acceptText = "Accept", declineText = "Decline" }) => {
  const { mutate } = useSWRConfig()

  const respond = (accept: boolean) =>
    respondToInvitation(invitation.id, accept)
      .then(() => {
        // group invitation, refetch group
        switch (invitation.type) {
          case 1:
          case 2:
            mutate(`/api/groups/${invitation.associated_id}`)
            break
          case 3:
            mutate(`/api/events/${invitation.associated_id}`)
            break
        }
      })
      .catch(console.error)
      .finally(() => mutate(`/api/invitations`))

  return (
    <div className={"flex gap-2"}>
      <button className={"btn btn-sm btn-success"} onClick={() => respond(true)}>
        {acceptText}
      </button>
      <button className={"btn btn-sm btn-error"} onClick={() => respond(false)}>
        {declineText}
      </button>
    </div>
  )
}
