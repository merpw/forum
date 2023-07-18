import { FC } from "react"

import { useUser } from "@/api/users/hooks"

const UserName: FC<{ userId: number }> = ({ userId }) => {
  const { user } = useUser(userId)

  if (!user) {
    return null
  }

  return <>{user.username}</>
}

export default UserName
