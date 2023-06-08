"use client"

import Link from "next/link"

import { logOut, useMe } from "@/api/auth"

const UserInfo = () => {
  const { user, isError, isLoading, isLoggedIn, mutate } = useMe()
  if (isError || isLoading) {
    return null
  }

  if (!isLoggedIn) {
    return (
      <Link className={"clickable flex font-light space-x-0.5"} href={"/login"}>
        <span>Login</span>
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={2.5}
          stroke={"currentColor"}
          className={"w-6 h-6 text-primary"}
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
    <div className={"flex gap-1 my-auto"}>
      <div className={"dropdown dropdown-end"}>
        <div tabIndex={0} className={"avatar btn btn-ghost btn-circle"}>
          <div className={"w-36 sm:w-48 rounded-full ring-2 ring-accent"}>
            {/* TODO: add Online/Offline ring. Online: ring-accent; Offline: ring-neutral */}
            <svg
              xmlns={"http://www.w3.org/2000/svg"}
              viewBox={"0 0 24 24"}
              fill={"currentColor"}
              className={"opacity-30"}
            >
              <path
                fillRule={"evenodd"}
                d={
                  "M18.685 19.097A9.723 9.723 0 0021.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 003.065 7.097A9.716 9.716 0 0012 21.75a9.716 9.716 0 006.685-2.653zm-12.54-1.285A7.486 7.486 0 0112 15a7.486 7.486 0 015.855 2.812A8.224 8.224 0 0112 20.25a8.224 8.224 0 01-5.855-2.438zM15.75 9a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
                }
                clipRule={"evenodd"}
              />
            </svg>
            {/* TODO: add Avatar*/}
          </div>
        </div>

        <ul
          tabIndex={0}
          className={
            "mt-3 p-2 shadow z-50 menu menu-compact dropdown-content bg-base-100 rounded-box w-52"
          }
        >
          <li className={"menu-title inline"}>
            <span className={"font-light"}>Hello,</span> {user?.name}
          </li>
          <li>
            <Link href={"/me"}>Profile</Link>
          </li>
          <li>
            <button
              type={"submit"}
              title={"Logout"}
              onClick={() =>
                logOut()
                  .then(() => mutate())
                  .catch(null)
              }
            >
              Logout
            </button>
          </li>
        </ul>
      </div>
    </div>
  )
}

export default UserInfo
