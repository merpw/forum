import { FC } from "react"
import pluralize from "pluralize"
import { MdGroupAdd, MdGroupRemove } from "react-icons/md"

import { Group } from "@/api/groups/hooks"
import GroupAvatar from "@/components/groups/Avatar"
import { useMe } from "@/api/auth/hooks"

const GroupInfo: FC<{ group: Group }> = ({ group }) => {
  return (
    <>
      <div className={"hero"}>
        <div className={"hero-content min-w-0 px-0"}>
          <div
            className={
              "card flex-shrink-0 w-full shadow-lg gradient-light dark:gradient-dark px-1 sm:px-3 pb-5"
            }
          >
            <div className={"card-body sm:flex-row sm:gap-7"}>
              <GroupAvatar group={group} size={200} className={"w-24 sm:w-48 m-auto self-center"} />
              <div className={"flex flex-col gap-2"}>
                <div
                  className={
                    "flex items-center justify-between font-light mb-5 text-center sm:text-left text-neutral"
                  }
                >
                  <h1
                    className={
                      "flex-shrink-0 text-3xl sm:text-4xl text-primary font-Yesteryear mr-10"
                    }
                  >
                    {group.title}
                  </h1>
                  <span className={"badge badge-neutral"}>Group</span>
                </div>

                <div className={""}>{group.description}</div>
                <div className={"font-light"}></div>
                <div className={"flex gap-7 justify-between items-center mt-auto"}>
                  <span className={"flex-shrink-0 text-primary"}>
                    {group.member_count} {pluralize("member", group.member_count)}
                  </span>
                  <MembershipButton group={group} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

const MembershipButton: FC<{ group: Group }> = ({ group }) => {
  const { user } = useMe()
  if (!user || user.id === group.creator_id) return null

  // group.member_status = 1

  if (group.member_status === 0) {
    return (
      <button className={"button w-fit px-7"}>
        <MdGroupAdd size={20} />
        Join
      </button>
    )
  }

  return (
    <button className={"btn btn-ghost text-info btn-neutral btn-sm"}>
      <MdGroupRemove size={20} />
      {group.member_status === 1 ? "Leave" : "Cancel"}
    </button>
  )
}

export default GroupInfo
