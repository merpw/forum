"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"

import { logOut, useMe } from "@/api/auth"

const UserInfo = () => {
  const { isError, isLoading, isLoggedIn, user, mutate } = useMe()

  const router = useRouter()

  if (isError || isLoading) {
    return null
  }

  if (!isLoggedIn) {
    return (
      <Link className={"clickable flex"} href={"/login"}>
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
          <span className={"clickable font-bold"}>{user?.name}</span>
        </Link>
      </div>
      <button
        title={"Logout"}
        className={"clickable cursor-pointer m-auto"}
        onClick={() =>
          logOut()
            .then(() => {
              router.refresh()
              mutate()
            })
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

export default UserInfo
