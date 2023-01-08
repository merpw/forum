import moment from "moment"
import Link from "next/link"
import { FC } from "react"
import { Post } from "../../custom"
import { Category, CommentsCount, ReactionsButtons } from "./reactions"

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

        {/* TODO: change to description */}
        <p>{post.content}</p>
      </div>

      <div className={"border-t flex flex-wrap p-3 gap-y-2 "}>
        <ReactionsButtons post={post} />
        <CommentsCount post={post} />
        <Category post={post} />

        <span className={"ml-auto"}>
          <span
            suppressHydrationWarning
            title={moment(post.date).local().format("DD.MM.YYYY HH:mm:ss")}
          >
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
