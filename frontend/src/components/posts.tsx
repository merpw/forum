import Link from "next/link"
import { FC } from "react"
import { Post } from "../custom"

export const PostList: FC<{ posts: Post[] }> = ({ posts }) => {
  if (posts.length == 0) {
    return <div>There are no posts yet...</div>
  }
  return <div>{posts.map(PostCard)}</div>
}

const PostCard = (post: Post, key: number) => (
  <div
    key={key}
    className={"rounded-lg border opacity-90 hover:opacity-100 hover:shadow dark:shadow-white mb-4"}
  >
    <div className={"m-5"}>
      <div className={"mb-3"}>
        <h1 className={"text-2xl hover:opacity-50 w-fit"}>
          <Link href={`/post/${post.id}`}>{post.title}</Link>
        </h1>
        <hr className={"mt-2"} />
      </div>

      <p>{post.content}</p>
    </div>

    <div className={"border-t px-5 py-2 flex justify-between"}>
      <span>{new Date(Date.parse(post.date)).toLocaleString("fi")}</span>
      <span>
        {"by "}
        <Link href={`/user/${post.author.id}`}>
          <span className={"text-xl hover:opacity-50"}>{post.author.name}</span>
        </Link>
      </span>
    </div>
  </div>
)
