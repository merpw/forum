// TODO: improve styles

import { useContext } from "react"

import { ChatSectionCollapsedContext } from "@/components/layout"
import { useChatIds } from "@/api/chats/chats"
import ChatList from "@/components/chats/ChatList"

const ChatsSection = () => {
  const { chatIds } = useChatIds()

  const { isCollapsed, setIsCollapsed } = useContext(ChatSectionCollapsedContext)

  if (isCollapsed) {
    return (
      <button
        onClick={() => setIsCollapsed(false)}
        className={"absolute right-2 backdrop-blur-2xl rounded-2xl"}
      >
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={1.5}
          stroke={"currentColor"}
          className={"w-6 h-6 m-2"}
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
    <div className={"max-h-full relative overflow-auto pr-2 w-full sm:w-1/2 md:w-1/4"}>
      <div className={"flex"}>
        <h1 className={"text-2xl"}>Chats</h1>
        <button className={"ml-auto"} onClick={() => setIsCollapsed(true)}>
          <svg
            xmlns={"http://www.w3.org/2000/svg"}
            fill={"none"}
            viewBox={"0 0 24 24"}
            strokeWidth={1.5}
            stroke={"currentColor"}
            className={"w-6 h-6 m-2"}
          >
            <path
              strokeLinecap={"round"}
              strokeLinejoin={"round"}
              d={"M11.25 4.5l7.5 7.5-7.5 7.5m-6-15l7.5 7.5-7.5 7.5"}
            />
          </svg>
        </button>
      </div>
      {chatIds ? <ChatList chatIds={chatIds} /> : <div>loading...</div>}
    </div>
  )
}

export default ChatsSection
