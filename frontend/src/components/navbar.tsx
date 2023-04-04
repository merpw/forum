import { NextComponentType } from "next"
import Link from "next/link"
import { logOut, useMe } from "@/api/auth"

// TODO: fix hydration error

const Navbar: NextComponentType = () => {
  return (
    <div className={"border-b mb-5 pb-2"}>
      <nav className={"flex justify-between flex-wrap"}>
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
  const { isError, isLoading, isLoggedIn, user, mutate } = useMe()
  if (isError || isLoading) {
    return null
  }

  if (!isLoggedIn) {
    return (
      <Link className={"hover:opacity-50 flex"} href={"/login"}>
        <span>Login</span>
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={1.5}
          stroke={"currentColor"}
          className={"w-6 h-6"}
        >
          <title>Login</title>
          <path
            strokeLinecap={"round"}
            strokeLinejoin={"round"}
            d={
              "M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75"
            }
          />
        </svg>
      </Link>
    )
  }
  return (
    <div className={"flex gap-1"}>
      <div className={"text-lg"}>
        {"Hello, "}
        <Link href={"/me"}>
          <span className={"font-bold hover:opacity-50"}>{user?.name}</span>
        </Link>
      </div>
      <button
        className={"cursor-pointer hover:opacity-50 m-auto"}
        onClick={() =>
          logOut()
            .then(() => mutate())
            .catch(null)
        }
      >
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={1.5}
          stroke={"currentColor"}
          className={"w-6 h-6"}
        >
          <title>Logout</title>
          <path
            strokeLinecap={"round"}
            strokeLinejoin={"round"}
            d={
              "M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15m3 0l3-3m0 0l-3-3m3 3H9"
            }
          />
        </svg>
      </button>
    </div>
  )
}

export default Navbar
