import { FC } from "react"
import { useSWRConfig } from "swr"

import respondToInvitation from "@/api/invitations/respond"
import { Invitation } from "@/api/invitations/hooks"

export const ResponseButtons: FC<{ invitation: Invitation }> = ({ invitation }) => {
  const { mutate } = useSWRConfig()

  const respond = (accept: boolean) =>
    respondToInvitation(invitation.id, accept)
      .then(() => {
        // group invitation, refetch group
        if (invitation.type === 1) mutate(`/api/groups/${invitation.associated_id}`)
      })
      .catch(console.error)
      .finally(() => mutate(`/api/invitations`))

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
