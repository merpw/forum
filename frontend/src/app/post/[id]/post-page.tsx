import { FC } from "react"

import { Comment, Post } from "@/custom"
import { Category, ReactionsButtons } from "@/components/posts/reactions"
import CommentForm from "@/app/post/[id]/CommentForm"
import Comments from "@/app/post/[id]/Comments"
import AutoDate from "@/components/AutoDate"
import UserLink from "@/components/UserLink"

import "highlight.js/styles/github-dark.css"

const PostPage: FC<{ post: Post; comments: Comment[] }> = ({ post, comments }) => {
  return (
    <div className={"max-w-3xl mx-auto"}>
      <div className={"rounded-lg border border-primary shadow-base-100 shadow-sm my-7"}>
        <div className={"mx-5"}>
          <h1 className={"start-dot end-dot flex flex-wrap text-2xl my-3"}>{post.title}</h1>
          <div
            className={"prose dark:prose-invert m-3 mb-5 text-lg font-light"}
            dangerouslySetInnerHTML={{ __html: post.content }}
          />
        </div>

        <div className={"bg-base-100 m-1 rounded-lg flex flex-wrap py-2 px-5 gap-x-1"}>
          <span className={"self-center mr-3"}>
            <ReactionsButtons post={post} />
          </span>
          <Category post={post} />
          <span className={"ml-auto font-light"}>
            <span className={"text-info"}>
              <AutoDate date={post.date} />
              {" by "}
            </span>
            <UserLink
              userId={post.author.id}
              className={"clickable text-lg font-Alatsi self-center"}
            >
              {post.author.username}
            </UserLink>
          </span>
        </div>
      </div>
      <div>
        <CommentForm postId={post.id} />
        <div className={"text-center mt-12"}>
          <h2 className={"tab tab-bordered tab-active cursor-default self-center mb-3"}>
            Comments
          </h2>
        </div>
        <Comments postId={post.id} fallback={comments} />
      </div>
    </div>
  )
}

export default PostPage
