import { FC } from "react"
import Link from "next/link"
import pluralize from "pluralize"
import { MdGroupAdd, MdGroupRemove, MdOutlineGroupAdd } from "react-icons/md"
import { HiChatAlt2 } from "react-icons/hi"

import { Group, useGroup, useGroups } from "@/api/groups/hooks"

const GroupList = () => {
  const { groups } = useGroups()
  if (!groups) return null

  return (
    <div className={"flex flex-col gap-2"}>
      <div className={"m-3"}>
        <Link href={"/group/create"} className={"button"}>
          <span className={"my-auto"}>
            <MdOutlineGroupAdd size={20} />
          </span>
          <span className={"ml-1 text-xs"}>Create a new group</span>
        </Link>
      </div>
      <span className={"text-info text-sm my-1"}>
        Total: {groups.length} {pluralize("group", groups.length)}
      </span>
      {groups.map((groupId) => (
        <GroupCard groupId={groupId} key={groupId} />
      ))}
    </div>
  )
}

const GroupCard: FC<{ groupId: number }> = ({ groupId }) => {
  const { group } = useGroup(groupId)

  if (!group) return null

  return (
    <div
      className={
        "w-full hover:bg-neutral hover:saturate-150 p-2 rounded-2xl flex justify-between items-center"
      }
    >
      <Link href={`/group/${groupId}`} className={"grow"}>
        <h3 className={"text-primary"}>{group.title}</h3>
        <div className={"text-info"}>
          {group.member_count} {pluralize("member", group.member_count)}
        </div>
      </Link>
      <GroupCardButton group={group} />
    </div>
  )
}

const size = 35

const GroupCardButton: FC<{ group: Group }> = ({ group }) => {
  if (group.member_status === 0) {
    // not a member
    return <MdGroupAdd size={size} className={"text-primary mx-1.5"} />
  }

  if (group.member_status === 2) {
    // pending
    return <MdGroupRemove size={size} className={"text-info mx-1.5"} />
  }

  return (
    <Link href={`/chat/g${group.id}`} className={"btn btn-ghost rounded-2xl px-1.5"}>
      <HiChatAlt2 size={size} className={"text-primary"} />
    </Link>
  )
}

export default GroupList
