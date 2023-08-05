import { Invitation } from "@/api/invitations/hooks"

const DUMMY_INVITATIONS: Invitation[] = [
  {
    id: 1,
    type: 0,
    associatedId: 1,
    userId: 1,
    timestamp: "2021-08-01T00:00:00.000Z",
  },
  {
    id: 2,
    type: 0,
    associatedId: 1,
    userId: 2,
    timestamp: "2021-08-01T00:00:00.000Z",
  },
  {
    id: 3,
    type: 0,
    associatedId: 1,
    userId: 3,
    timestamp: "2021-08-01T00:00:00.000Z",
  },
]

export const removeDummyInvitation = (id: number) => {
  DUMMY_INVITATIONS.splice(
    DUMMY_INVITATIONS.findIndex((invitation) => invitation.id === id),
    1
  )
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const dummyFetcher = async (url: string): Promise<any> => {
  if (url === "/api/invitations") {
    return Promise.resolve(DUMMY_INVITATIONS.map((invitation) => invitation.id))
  }
  if (url.match(/\/api\/invitations\/\d+/)) {
    const id = url.split("/").pop()
    return Promise.resolve(DUMMY_INVITATIONS.find((invitation) => invitation.id === Number(id)))
  }
  return Promise.reject(new Error("Invalid URL"))
}
