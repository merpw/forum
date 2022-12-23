import { NextPage } from "next/types"
import { useUser } from "../fetch/user"

const UserPage: NextPage = () => {
  const { user, isError, isLoading, isLoggedIn } = useUser()
  if (isError) {
    return <div>Error</div>
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (!isLoggedIn) {
    return <div>Not logged in</div>
  }

  return (
    <div>
      <h1 className={"text-xl"}>{user?.name}</h1>
    </div>
  )
}

export default UserPage
