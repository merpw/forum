"use client"

import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useEffect,
  useRef,
  useState,
} from "react"
import { Provider } from "react-redux"
import Link from "next/link"

import { useMe } from "@/api/auth/hooks"
import Navbar from "@/components/navbar"
import store from "@/store/store"
import ChatsSection from "@/components/ChatSection"

export const ChatSectionCollapsedContext = createContext({
  isCollapsed: false,
  setIsCollapsed: (() => void 0) as Dispatch<SetStateAction<boolean>>,
})

export default function Layout({ children }: { children: ReactNode }) {
  const { user } = useMe()
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // TODO: add scroll restoration
    if (scrollRef.current) {
      scrollRef.current.scrollTo(0, 0)
    }
  }, [children])

  const [isCollapsed, setIsCollapsed] = useState(true)

  useEffect(() => {
    if (window.innerWidth >= 640) {
      setIsCollapsed(false)
    }
  }, [])

  return (
    <Provider store={store}>
      <div className={"container m-auto p-5 min-h-screen break-words"}>
        <header>
          <Navbar />
        </header>
        <div className={"mx-5 flex gap-3 flex-row h-[calc(100vh-10.5rem)] relative"}>
          <main
            ref={scrollRef}
            className={
              "max-h-full overflow-auto grow" +
              " " +
              (user && !isCollapsed ? "hidden sm:block w-3/4" : "")
            }
          >
            {children}
          </main>
          {user && (
            <ChatSectionCollapsedContext.Provider value={{ isCollapsed, setIsCollapsed }}>
              <ChatsSection />
            </ChatSectionCollapsedContext.Provider>
          )}
        </div>
        <footer>
          <div className={"text-center font-light text-neutral border-t border-base-100 pt-3 mt-3"}>
            <p>
              {"Made with "}
              <Link href={"https://github.com/merpw/forum"} className={"gradient-text clickable"}>
                {"❤"}
              </Link>
              {" by "}
              <span className={"font-light"}>{"basket of kittens"}</span>
            </p>
          </div>
        </footer>
      </div>
    </Provider>
  )
}
