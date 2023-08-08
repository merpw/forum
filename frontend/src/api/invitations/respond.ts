import axios from "axios"

const respond = async (invitationId: number, response: boolean) =>
  axios.post(`/api/invitations/${invitationId}/respond`, { response })

export default respond
