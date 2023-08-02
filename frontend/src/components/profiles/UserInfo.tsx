/* eslint-disable import/no-named-as-default-member */
/* https://github.com/iamkun/dayjs/issues/1242 */
import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import { FC } from "react"

import { User } from "@/custom"

/* "2000-01-24" -> "23 years"
 * "2021-01-24" -> "baby 👶" */
const calculateAge = (dob: string): string | null => {
  dayjs.extend(relativeTime)

  const parsedDob = dayjs(dob, "YYYY-MM-DD")
  if (!parsedDob.isValid()) return null

  const age = parsedDob.fromNow(true)
  return age.includes("year") ? age + " old" : "baby 👶"
}
export const UserInfo: FC<{ user: User }> = ({ user }) => {
  const age = user.dob ? calculateAge(user.dob) : null
  return (
    <>
      <div
        className={"flex flex-col self-center font-light mb-5 text-center sm:text-left text-info"}
      >
        {"Hey, "}
        <span className={"text-3xl sm:text-4xl text-primary font-Yesteryear mx-1"}>
          {user.username}
        </span>
        {"Forgot who you are?"}
      </div>
      <div className={"font-light"}>
        {user.first_name || user.last_name ? (
          <p>
            <span className={"font-light text-info"}>{"Full name"}</span>
            <span
              className={"font-normal start-dot"}
            >{`${user.first_name} ${user.last_name}`}</span>
          </p>
        ) : null}
        {age ? (
          <p>
            <span className={"font-light text-info"}>{"Age"}</span>
            <span className={"font-normal start-dot"}>{age}</span>
          </p>
        ) : null}
        {user.gender ? (
          <p>
            <span className={"font-light text-info"}>{"Gender"}</span>
            <span className={"font-normal start-dot"}>{user.gender}</span>
          </p>
        ) : null}
        <p>
          <span className={"font-light text-info"}>{"Email"}</span>
          <span className={"font-normal start-dot"}>{user.email}</span>
        </p>
      </div>
    </>
  )
}
