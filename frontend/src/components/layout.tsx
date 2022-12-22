import { ReactNode } from "react"
import Navbar from "./navbar"

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className={"container mx-auto p-5 min-h-screen"}>
      <header>
        <Navbar />
      </header>
      <main>{children}</main>
      <hr className={"my-7 opacity-50"} />
    </div>
  )
}
