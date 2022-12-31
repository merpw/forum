import moment from "moment"
import Link from "next/link"
import { FC } from "react"
import { Post } from "../../custom"
import { ReactionsButtons } from "./reactions"

export const PostList: FC<{ posts: Post[] }> = ({ posts }) => {
  if (posts.length == 0) {
    return <div>There are no posts yet...</div>
  }
  return <div>{posts.map(PostCard)}</div>
}

const PostCard = (post: Post, key: number) => {
  return (
    <div
      key={key}
      className={
        "rounded-lg border opacity-90 hover:opacity-100 hover:shadow dark:shadow-white mb-4"
      }
    >
      <div className={"m-5"}>
        <div className={"mb-3"}>
          <h1 className={"text-2xl hover:opacity-50 max-w-fit"}>
            <Link href={`/post/${post.id}`}>{post.title}</Link>
          </h1>
          <hr className={"mt-2"} />
        </div>

        <p>{post.content}</p>
      </div>

      <div className={"border-t flex flex-wrap p-3 gap-y-2 "}>
        <ReactionsButtons post={post} />
        <Link href={`/category/${post.category}`} className={"hover:opacity-50 flex"}>
          <span className={"text-xl capitalize"}>{post.category}</span>
          <span className={"pt-1 ml-1"}>
            <svg
              xmlns={"http://www.w3.org/2000/svg"}
              fill={"none"}
              viewBox={"0 0 24 24"}
              strokeWidth={1.5}
              stroke={"currentColor"}
              className={"w-6 h-6"}
            >
              <path
                strokeLinecap={"round"}
                strokeLinejoin={"round"}
                d={
                  "M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3z"
                }
              />
              <path strokeLinecap={"round"} strokeLinejoin={"round"} d={"M6 6h.008v.008H6V6z"} />
            </svg>
          </span>
        </Link>

        <span className={"ml-auto"}>
          <span title={moment(post.date).local().format("DD.MM.YYYY HH:mm:ss")}>
            {moment(post.date).fromNow()}
          </span>
          {" by "}
          <span className={"text-xl hover:opacity-50"}>
            <Link href={`/user/${post.author.id}`}>{post.author.name}</Link>
          </span>
        </span>
      </div>
    </div>
  )
}
