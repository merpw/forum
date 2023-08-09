"use client"

import { FC } from "react"
import { SWRConfig } from "swr"

import { Group, useGroup } from "@/api/groups/hooks"
import GroupInfo from "@/components/groups/GroupInfo"
import { PostList } from "@/components/posts/list"
import { useGroupPosts } from "@/api/groups/posts"

const GroupPage: FC<{ groupId: number }> = ({ groupId }) => {
  const { group } = useGroup(groupId)

  if (!group) return null

  return (
    <>
      <GroupInfo group={group} />
      <div className={"mt-5"}>
        <div className={"text-center"}>
          <h2 className={"tab tab-bordered tab-active cursor-default self-center mb-3"}>Posts</h2>
        </div>
        {group.member_status === 1 ? (
          <GroupPosts groupId={groupId} />
        ) : (
          <div className={"text-info text-center mt-5 mb-7"}>
            You need to be a member to see posts
          </div>
        )}
      </div>
    </>
  )
}

const GroupPosts: FC<{ groupId: number }> = ({ groupId }) => {
  const { posts } = useGroupPosts(groupId)
  if (!posts) return null

  return <PostList posts={posts} />
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
