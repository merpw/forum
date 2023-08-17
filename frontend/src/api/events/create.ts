import axios from "axios"

const createEvent = async (
  event: { title: string; description: string; time_and_date: string },
  groupId: number
) =>
  axios
    .post(`/api/groups/${groupId}/events/create`, event)
    .then((res) => res.data)
    .catch((err) => {
      throw new Error(err.response.data)
    })

export default createEvent
