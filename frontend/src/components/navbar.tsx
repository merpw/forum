import { NextComponentType } from "next"
import Link from "next/link"
import { logOut, useUser } from "../fetch/user"

const Navbar: NextComponentType = () => {
  return (
    <div className={"border-b mb-5 pb-2"}>
      <nav className={"flex justify-between"}>
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
  const { isError, isLoading, isLoggedIn, user, mutate } = useUser()

  if (isError || isLoading) {
    return null
  }

  if (!isLoggedIn) {
    return (
      <Link className={"hover:opacity-50"} href={"/login"}>
        Login
      </Link>
    )
  }
  return (
    <div className={"flex gap-2"}>
      <div className={"text-lg"}>
        {"Hello, "}
        <Link href={"/me"}>
          <span className={"font-bold hover:opacity-50"}>{user?.name}</span>
        </Link>
      </div>
      <div>
        <span
          className={"cursor-pointer hover:opacity-50"}
          onClick={() => logOut().then(() => mutate())}
        >
          Logout
        </span>
      </div>
    </div>
  )
}

export default Navbar
