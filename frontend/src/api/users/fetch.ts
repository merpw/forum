import { User } from "../../custom"
import { users } from "../dummy"

export const getUsersLocal = (): Promise<User[]> => Promise.resolve(users)
export const getUserLocal = (id: number): Promise<User | undefined> =>
  Promise.resolve(users.find((user) => user.id == id))
