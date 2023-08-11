import { FC } from "react"
import { useSWRConfig } from "swr"

import respondToInvitation from "@/api/invitations/respond"

export const ResponseButtons: FC<{ invitationId: number }> = ({ invitationId }) => {
  const { mutate } = useSWRConfig()

  const respond = (accept: boolean) =>
    respondToInvitation(invitationId, accept)
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
