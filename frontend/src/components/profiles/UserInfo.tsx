/* eslint-disable import/no-named-as-default-member */
/* https://github.com/iamkun/dayjs/issues/1242 */
import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import { FC, useState } from "react"
import pluralize from "pluralize"

import { User } from "@/custom"
import Avatar from "@/components/Avatar"
import { togglePrivacy } from "@/api/users/hooks"

export const UserInfo: FC<{ user: User; isOwnProfile?: boolean }> = ({ user, isOwnProfile }) => {
  const age = user.dob ? calculateAge(user.dob) : null
  return (
    <>
      <div className={"hero"}>
        <div className={"hero-content px-0"}>
          <div
            className={
              "card flex-shrink-0 w-full shadow-lg gradient-light dark:gradient-dark px-1 sm:px-3 pb-5"
            }
          >
            <div className={"card-body sm:flex-row sm:gap-5"}>
              <Avatar user={user} size={200} className={"w-24 sm:w-48 m-auto self-center"} />
              <div className={"self-center text-sm"}>
                <div
                  className={
                    "flex flex-col self-center font-light mb-5 text-center sm:text-left text-info"
                  }
                >
                  {isOwnProfile ? "Hey, " : "Hey! I'm "}
                  <span className={"text-3xl sm:text-4xl text-primary font-Yesteryear mx-1"}>
                    {user.username}
                  </span>
                  {isOwnProfile ? (
                    <>
                      Forgot who you are?
                      <PrivacyToggle user={user} />
                    </>
                  ) : (
                    user.privacy && (
                      <span className={"badge badge-outline"}>
                        <svg
                          xmlns={"http://www.w3.org/2000/svg"}
                          fill={"none"}
                          viewBox={"0 0 24 24"}
                          strokeWidth={1.5}
                          stroke={"currentColor"}
                          className={"w-3.5 h-3.5 mr-1"}
                        >
                          <path
                            strokeLinecap={"round"}
                            strokeLinejoin={"round"}
                            d={
                              "M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z"
                            }
                          />
                        </svg>
                        private
                      </span>
                    )
                  )}
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
                  {user.email ? (
                    <p>
                      <span className={"font-light text-info"}>{"Email"}</span>
                      <span className={"font-normal start-dot"}>{user.email}</span>
                    </p>
                  ) : null}
                </div>
              </div>
            </div>
            <div className={"flex gap-2 justify-center"}>
              <span className={"text-primary"}>
                {user.followers_count} {pluralize("follower", user.followers_count)}
              </span>
              <span className={"text-info"}>•</span>
              <span className={"text-secondary"}>{user.following_count} following</span>
            </div>

            {user.bio && (
              <div className={"mt-2 text-center"}>
                <div className={"font-light text-info start-dot end-dot mb-1"}>About me</div>
                <div className={"text-sm"}>{user.bio}</div>
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  )
}

/* "2000-01-24" -> "23 years"
 * "2021-01-24" -> "baby 👶" */
const calculateAge = (dob: string): string | null => {
  dayjs.extend(relativeTime)

  const parsedDob = dayjs(dob, "YYYY-MM-DD")
  if (!parsedDob.isValid()) return null

  const age = parsedDob.fromNow(true)
  return age.includes("year") ? age + " old" : "baby 👶"
}

const PrivacyToggle: FC<{ user: User }> = ({ user }) => {
  const [privacy, setPrivacy] = useState(user.privacy)
  return (
    <div className={"form-control mt-3"}>
      <label className={"label cursor-pointer"}>
        <span className={"label-text"}>Private</span>
        <input
          type={"checkbox"}
          className={"toggle ml-3 mr-auto"}
          checked={Boolean(privacy)}
          onChange={() => togglePrivacy().then(setPrivacy)}
        />
      </label>
    </div>
  )
}
