import { FC } from "react"
import Link from "next/link"
import pluralize from "pluralize"
import { MdGroupAdd, MdGroupRemove } from "react-icons/md"
import { HiChatAlt2 } from "react-icons/hi"

import { Group, useGroup, useGroups } from "@/api/groups/hooks"

const GroupList = () => {
  const { groups } = useGroups()
  if (!groups) return null

  return (
    <div className={"flex flex-col gap-2"}>
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
    <div className={"w-full hover:bg-neutral hover:saturate-150 p-2 rounded-2xl"}>
      <Link href={`/group/${groupId}`} className={"flex justify-between items-center"}>
        <div>
          <h3 className={"text-primary"}>{group.title}</h3>
          <div className={"text-info"}>
            {group.member_count} {pluralize("member", group.member_count)}
          </div>
        </div>
        <GroupCardButton group={group} />
      </Link>
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
