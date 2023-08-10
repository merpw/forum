"use client"

import { FC } from "react"
import { SWRConfig } from "swr"
import Link from "next/link"

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
      {group.member_status === 1 ? (
        <div className={"m-3 text-center"}>
          <Link href={`/group/${groupId}/compose`} className={"button"}>
            <span className={"my-auto"}>
              <svg
                xmlns={"http://www.w3.org/2000/svg"}
                fill={"none"}
                viewBox={"0 0 24 24"}
                stroke={"currentColor"}
                className={"w-5 h-5"}
              >
                <path
                  strokeLinecap={"round"}
                  strokeLinejoin={"round"}
                  strokeWidth={1.5}
                  d={
                    "M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
                  }
                />
              </svg>
            </span>
            <span className={"ml-1 text-xs"}>Create a new post</span>
          </Link>
        </div>
      ) : null}
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
