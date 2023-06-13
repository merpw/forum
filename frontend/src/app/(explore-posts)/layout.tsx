import { ReactNode } from "react"
import Link from "next/link"

import CategoryButtons from "@/components/CategoryButtons"
import { getCategoriesLocal } from "@/api/posts/edge"

const Layout = async ({ children }: { children: ReactNode }) => {
  const categories = await getCategoriesLocal()

  return (
    <>
      <div className={"flex flex-wrap justify-center flex-col text-center"}>
        <div className={"my-2 font-Yesteryear text-3xl text-primary opacity-30 dark:opacity-50"}>
          <span className={"text-neutral-content"}> Welcome aboard!</span>
          {/* TODO: add User's name */}
        </div>
        <CreatePostButton />
        <CategoryButtons categories={categories} />
      </div>
      {children}
    </>
  )
}

const CreatePostButton = () => {
  return (
    <div className={"m-3"}>
      <Link href={"/create"} className={"button"}>
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
  )
}

export default Layout
