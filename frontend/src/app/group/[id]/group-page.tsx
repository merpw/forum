"use client"

import { FC, useRef } from "react"
import { SWRConfig } from "swr"
import Link from "next/link"
import { HiOutlineUserAdd } from "react-icons/hi"
import { HiOutlinePencilSquare } from "react-icons/hi2"

import { Group, useGroup } from "@/api/groups/hooks"
import GroupInfo from "@/components/groups/GroupInfo"
import { PostList } from "@/components/posts/list"
import { useGroupPosts } from "@/api/groups/posts"
import { usePostList } from "@/api/posts/hooks"
import InviteFollowersForm from "@/components/forms/groups/InviteFollowersForm"

const GroupPage: FC<{ groupId: number }> = ({ groupId }) => {
  const { group } = useGroup(groupId)

  const invitePopupRef = useRef<HTMLDialogElement>(null)

  if (!group) return null

  return (
    <>
      <GroupInfo group={group} />
      {group.member_status === 1 ? (
        <div className={"flex justify-center gap-3"}>
          <button className={"button"} onClick={() => invitePopupRef.current?.showModal()}>
            <HiOutlineUserAdd className={"w-5 h-5"} />
            Invite
          </button>
          <dialog ref={invitePopupRef} className={"modal"}>
            <InviteFollowersForm groupId={group.id} />
            <form method={"dialog"} className={"modal-backdrop"}>
              <button>close</button>
            </form>
          </dialog>
          <Link href={`/group/${groupId}/compose`} className={"button"}>
            <HiOutlinePencilSquare className={"w-5 h-5"} />
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
  const { posts: postIds } = useGroupPosts(groupId)

  const { posts } = usePostList(postIds || [])

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
