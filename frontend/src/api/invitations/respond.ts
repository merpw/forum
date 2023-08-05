import { removeDummyInvitation } from "@/api/invitations/dummy"

const respond = async (invitationId: number, response: boolean) => {
  removeDummyInvitation(invitationId)
}

export default respond
