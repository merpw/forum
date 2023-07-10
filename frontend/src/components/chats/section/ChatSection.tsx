// TODO: improve styles

import { useContext, useEffect, useState } from "react"

import { ChatSectionCollapsedContext } from "@/components/layout"
import { useChatIds } from "@/api/chats/chats"
import ChatList from "@/components/chats/section/ChatList"
import UserList from "@/components/chats/section/UserList"

const ChatsSection = () => {
  const { chatIds } = useChatIds()

  const [tab, setTab] = useState<"chats" | "users">("chats")

  useEffect(() => {
    if (chatIds && chatIds.length === 0) {
      setTab("users")
    }
  }, [chatIds])

  const { isCollapsed, setIsCollapsed } = useContext(ChatSectionCollapsedContext)

  if (isCollapsed) {
    return (
      <button
        onClick={() => setIsCollapsed(false)}
        className={"absolute right-2 backdrop-blur-2xl rounded-2xl p-3"}
      >
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={3}
          stroke={"currentColor"}
          className={"w-6 h-6 m-2 text-primary"}
        >
          <path
            strokeLinecap={"round"}
            strokeLinejoin={"round"}
            d={"M18.75 19.5l-7.5-7.5 7.5-7.5m-6 15L5.25 12l7.5-7.5"}
          />
        </svg>
      </button>
    )
  }

  return (
    <div className={"bg-base-100 max-h-full relative overflow-auto p-3 w-full sm:w-1/2 md:w-1/4"}>
      <div className={"flex flex-column mb-2"}>
        <div className={"space-y-3"}>
          <div className={""}>
            {/* TODO: add User list */}
            <ul className={"tab tab-md p-0 font-bold"}>
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
          </div>
        </div>

        <button className={"ml-auto"} onClick={() => setIsCollapsed(true)}>
          <svg
            xmlns={"http://www.w3.org/2000/svg"}
            fill={"none"}
            viewBox={"0 0 24 24"}
            strokeWidth={3}
            stroke={"currentColor"}
            className={"w-6 h-6 m-2 text-primary"}
          >
            <path
              strokeLinecap={"round"}
              strokeLinejoin={"round"}
              d={"M11.25 4.5l7.5 7.5-7.5 7.5m-6-15l7.5 7.5-7.5 7.5"}
            />
          </svg>
        </button>
      </div>
      {tab === "chats" && chatIds && <ChatList chatIds={chatIds} />}
      {tab === "users" && <UserList />}
    </div>
  )
}

export default ChatsSection
