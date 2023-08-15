"use client"

import { FC, useState } from "react"
import { SWRConfig } from "swr"
import Link from "next/link"
import { HiChatAlt2, HiOutlineUserAdd } from "react-icons/hi"
import { HiOutlineCalendarDays, HiOutlinePencilSquare } from "react-icons/hi2"

import { Group, useGroup } from "@/api/groups/hooks"
import GroupInfo from "@/components/groups/GroupInfo"
import { PostList } from "@/components/posts/list"
import { useGroupPosts } from "@/api/groups/posts"
import { usePostList } from "@/api/posts/hooks"
import InviteFollowersForm from "@/components/forms/groups/InviteFollowersForm"
import EventList from "@/components/events/EventList"
import { useGroupEvents } from "@/api/events/hooks"
import CreateEventForm from "@/components/forms/groups/CreateEventForm"

const GroupPage: FC<{ groupId: number }> = ({ groupId }) => {
  const { group } = useGroup(groupId)

  if (!group) return null

  return (
    <>
      <GroupInfo group={group} />
      {group.member_status === 1 ? (
        <div className={"flex justify-center gap-3"}>
          <button
            className={"button h-12 w-12 rounded-full"}
            onClick={(e) => (e.currentTarget.nextElementSibling as HTMLDialogElement).showModal()}
            title={"Invite followers to the group"}
          >
            <HiOutlineUserAdd className={"w-full h-full"} />
          </button>
          <dialog className={"modal"}>
            <InviteFollowersForm groupId={group.id} />
            <form method={"dialog"} className={"modal-backdrop"}>
              <button>close</button>
            </form>
          </dialog>
          <Link
            href={`/group/${groupId}/compose`}
            className={"button h-12 w-12 rounded-full"}
            title={"Create a new post in the group"}
          >
            <HiOutlinePencilSquare className={"w-full h-full"} />
          </Link>
          <Link
            href={`/chat/g${groupId}`}
            className={"button h-12 w-12 rounded-full"}
            title={"Open the group chat"}
          >
            <HiChatAlt2 className={"w-full h-full"} />
          </Link>
          <button
            className={"button h-12 w-12 rounded-full"}
            onClick={(e) => (e.currentTarget.nextElementSibling as HTMLDialogElement).showModal()}
            title={"Create event"}
          >
            <HiOutlineCalendarDays className={"w-full h-full"} />
          </button>
          <dialog className={"modal"}>
            <CreateEventForm groupId={group.id} />
            <form method={"dialog"} className={"modal-backdrop"}>
              <button>close</button>
            </form>
          </dialog>
        </div>
      ) : null}

      <GroupFooter groupId={groupId} />
    </>
  )
}

const GroupFooter: FC<{ groupId: number }> = ({ groupId }) => {
  const [activeTab, setActiveTab] = useState<number>(0)

  const ActiveTabComponent = Tabs[activeTab].Component

  const { group } = useGroup(groupId)

  if (!group || group.member_status === undefined) return null

  if (group.member_status === 0)
    return (
      <div className={"text-info text-center mt-5 mb-7"}>
        You need to be a member to see posts and events
      </div>
    )

  return (
    <div className={"mt-5"}>
      <div className={"text-center mb-2"}>
        {Tabs.map(({ name }, key) => (
          <h2
            key={key}
            className={`tab tab-bordered ${activeTab === key ? "tab-active" : ""}`}
            onClick={() => setActiveTab(key)}
          >
            {name}
          </h2>
        ))}
      </div>
      <ActiveTabComponent groupId={groupId} />
    </div>
  )
}

const GroupPosts: FC<{ groupId: number }> = ({ groupId }) => {
  const { posts: postIds } = useGroupPosts(groupId)

  const { posts } = usePostList(postIds || [])

  if (!posts) return null

  return <PostList posts={posts} />
}

const GroupEvents: FC<{ groupId: number }> = ({ groupId }) => {
  const { events } = useGroupEvents(groupId)

  if (!events) return null

  return <EventList events={events} />
}

const Tabs = [
  {
    name: "Posts",
    Component: GroupPosts,
  },
  {
    name: "Events",
    Component: GroupEvents,
  },
] as const

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
