import axios from "axios"

export const joinEvent = async (eventId: number) =>
  axios.post(`/api/events/${eventId}/going`).then((res) => res.data)

export const leaveEvent = async (eventId: number) =>
  axios.post(`/api/events/${eventId}/leave`).then((res) => res.data)
