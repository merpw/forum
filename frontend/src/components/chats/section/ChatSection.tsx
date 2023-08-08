// TODO: improve styles

import { FC, useCallback, useContext, useEffect, useState } from "react"
import { usePathname } from "next/navigation"

import { ChatSectionCollapsedContext } from "@/components/layouts/Layout"
import { useChatIds } from "@/api/chats/chats"
import ChatList from "@/components/chats/section/ChatList"
import UserList from "@/components/UserList"
import { useAppSelector } from "@/store/hooks"

const ChatsSection = () => {
  const { chatIds } = useChatIds()

  const [tab, setTab] = useState<"chats" | "users">("chats")

  useEffect(() => {
    if (chatIds && chatIds.length === 0) {
      setTab("users")
    }
  }, [chatIds])

  const { isCollapsed } = useContext(ChatSectionCollapsedContext)

  const pathname = usePathname()
  const collapseIfMobile = useCollapseIfMobile()

  useEffect(() => {
    collapseIfMobile()
  }, [collapseIfMobile, pathname])

  return (
    <>
      <CollapseButton />
      {!isCollapsed && (
        <div
          className={"bg-base-100 overscroll-contain overflow-auto p-3 w-full md:w-7/12 lg:w-1/3"}
        >
          <ul className={"tab tab-md p-0 font-bold mb-5"}>
            <li>
              <button
                type={"button"}
                className={
                  "tab tab-bordered space-y-5" + " " + (tab === "chats" ? "tab-active" : "")
                }
                onClick={() => setTab("chats")}
              >
                Chats
              </button>
            </li>
            <li>
              <button
                type={"button"}
                className={
                  "tab tab-bordered space-y-5" + " " + (tab === "users" ? "tab-active" : "")
                }
                onClick={() => setTab("users")}
              >
                Users
              </button>
            </li>
          </ul>
          {tab === "chats" && chatIds && <ChatList chatIds={chatIds} />}
          {tab === "users" && <UserList />}
        </div>
      )}
    </>
  )
}

const CollapseButton: FC = () => {
  const { setIsCollapsed } = useContext(ChatSectionCollapsedContext)

  const unreadCount = useAppSelector((state) => state.chats.unreadMessagesChatIds.length)

  return (
    <button
      onClick={() => setIsCollapsed((prev) => !prev)}
      className={
        "absolute hover:text-opacity-70 right-5 top-3 rounded-2xl p-1 text-primary z-40 backdrop-blur"
      }
    >
      <div className={"relative w-16 h-12"}>
        <span className={"absolute left-0 top-0"}>
          <svg
            xmlns={"http://www.w3.org/2000/svg"}
            fill={"none"}
            viewBox={"0 0 24 24"}
            strokeWidth={2}
            stroke={"currentColor"}
            className={"w-8 h-8 scale-x-[-1]"}
          >
            <path
              strokeLinecap={"round"}
              strokeLinejoin={"round"}
              d={
                "M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 01-2.555-.337A5.972 5.972 0 015.41 20.97a5.969 5.969 0 01-.474-.065 4.48 4.48 0 00.978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25z"
              }
            />
          </svg>
        </span>
        <span className={"absolute bottom-0 right-0"}>
          <svg
            xmlns={"http://www.w3.org/2000/svg"}
            fill={"none"}
            viewBox={"0 0 24 24"}
            strokeWidth={2}
            stroke={"currentColor"}
            className={"w-9 h-9"}
          >
            <path
              strokeLinecap={"round"}
              strokeLinejoin={"round"}
              d={
                "M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"
              }
            />
          </svg>
        </span>
        {unreadCount > 0 && (
          <span className={"absolute badge badge-secondary -right-3 rounded-full font-bold py-3"}>
            {unreadCount}
          </span>
        )}
      </div>
    </button>
  )
}

export const useCollapseIfMobile = () => {
  const { setIsCollapsed } = useContext(ChatSectionCollapsedContext)

  return useCallback(() => {
    if (window.innerWidth < 768) {
      setIsCollapsed(true)
    }
  }, [setIsCollapsed])
}

export default ChatsSection
