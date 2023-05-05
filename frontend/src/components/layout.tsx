import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react"
import { Provider } from "react-redux"

import { useMe } from "@/api/auth"
import ChatList from "@/components/chats/list"
import Navbar from "@/components/navbar"
import store from "@/store/store"
import { useChatIds } from "@/api/chats/chats"

export const ChatSectionCollapsedContext = createContext({
  isCollapsed: false,
  setIsCollapsed: (() => void 0) as Dispatch<SetStateAction<boolean>>,
})

export default function Layout({
  children,
  withChat = true,
}: {
  children: ReactNode
  withChat?: boolean
}) {
  const { user } = useMe()
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // TODO: add scroll restoration
    if (scrollRef.current) {
      scrollRef.current.scrollTo(0, 0)
    }
  }, [children])

  const [isCollapsed, setIsCollapsed] = useState(false)

  return (
    <Provider store={store}>
      <div className={"container max-w-screen-lg mx-auto break-words"}>
        <header className={"m-5"}>
          <Navbar />
        </header>
        <div className={"mx-5 flex gap-3 flex-row h-[calc(100vh-5.5rem)] relative"}>
          <main
            ref={scrollRef}
            className={
              "max-h-full overflow-auto grow" + " " + (!isCollapsed ? "hidden sm:block w-3/4" : "")
            }
          >
            {children}
          </main>
          {withChat && user !== null && (
            <ChatSectionCollapsedContext.Provider value={{ isCollapsed, setIsCollapsed }}>
              <ChatsSection />
            </ChatSectionCollapsedContext.Provider>
          )}
        </div>
      </div>
    </Provider>
  )
}

// TODO: improve styles, add mobile support

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
    <div
      className={
        "w-1/4 max-h-full relative overflow-auto pr-2" +
        " " +
        (!isCollapsed && "w-full sm:w-1/2 md:w-1/4")
      }
    >
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
