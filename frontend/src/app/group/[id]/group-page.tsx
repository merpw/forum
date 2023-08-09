"use client"

import { FC } from "react"
import { SWRConfig } from "swr"

import { Group, useGroup } from "@/api/groups/hooks"
import GroupInfo from "@/components/groups/GroupInfo"

const GroupPage: FC<{ groupId: number }> = ({ groupId }) => {
  const { group } = useGroup(groupId)

  if (!group) return null

  return <GroupInfo group={group} />
}

const GroupPageWrapper: FC<{ group: Group }> = ({ group }) => {
  return (
    <SWRConfig
      value={{
        fallback: {
          [`/api/groups/${group.id}`]: group,
        },
      }}
    >
      <GroupPage groupId={group.id} />
    </SWRConfig>
  )
}

export default GroupPageWrapper
