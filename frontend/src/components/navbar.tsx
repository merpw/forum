import { NextComponentType } from "next"
import Link from "next/link"
import { useUser } from "../fetch/user"

const Navbar: NextComponentType = () => {
  return (
    <div className={"border-b mb-5 pb-2"}>
      <nav className={"flex justify-between mx-10 mb-5S"}>
        <div>
          <Link href={"/"} className={"text-3xl hover:opacity-50"}>
            FORUM
          </Link>
        </div>
        <div className={"ml-auto"}>
          <UserInfo />
        </div>
      </nav>
    </div>
  )
}

const UserInfo = () => {
  const { isError, isLoading, isLoggedIn, user } = useUser()

  if (isError) {
    return null
  }

  if (isLoading) {
    // TODO: add placeholder
    return <div>Loading...</div>
  }
  if (!isLoggedIn) {
    return null
  }
  return (
    <div className={"text-lg"}>
      {"Hello, "}
      <Link href={"/me"}>
        <span className={"font-bold hover:opacity-50"}>{user?.name}</span>
      </Link>
    </div>
  )
}

export default Navbar
