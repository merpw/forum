import { GetStaticPaths, GetStaticProps, NextPage } from "next"
import { Post, Comment } from "../../custom"

import Link from "next/link"
import { FC } from "react"
import Head from "next/head"
import moment from "moment"
import {getPostLocal, getPostsLocal} from "../../api/posts/fetch";

const PostPage: NextPage<{ post: Post }> = ({ post }) => {
  return (
    <>
      <Head>
        <title>{`${post.title} - Forum`}</title>
      </Head>
      <div className={"m-5 "}>
        <div className={"mb-3"}>
          <h1 className={"text-3xl mb-2 "}>{post.title}</h1>
          <hr />
        </div>
        <p>{post.content}</p>
        <hr className={"mt-4"} />
        <div className={"border-t py-2 flex justify-between"}>
          <span title={moment(post.date).local().format("DD.MM.YYYY HH:mm:ss")}>
            {moment(post.date).fromNow()}
          </span>
          <span>
            {"by "}
            <Link href={`/user/${post.author.id}`}>
              <span className={"text-xl hover:opacity-50"}>{post.author.name}</span>
            </Link>
          </span>
        </div>
        <h2 className={"text-2xl my-4"}>Comments:</h2>
        <div>
          <Comments comments={post.comments} />
        </div>
      </div>
    </>
  )
}

const Comments: FC<{ comments: Comment[] }> = ({ comments }) => {
  if (comments.length == 0) {
    return <div>There are no comments yet, write one first!</div>
  }

  return (
    <>
      {comments.map((comment, key) => (
        <div className={"border m-2 rounded p-5"} key={key}>
          <Link href={`/user/${comment.author.id}`}>
            <h3 className={"text-lg hover:opacity-50"}>{comment.author.name}</h3>
          </Link>
          <p>{comment.text}</p>
          <span>{new Date(Date.parse(comment.date)).toLocaleString("fi")}</span>
        </div>
      ))}
    </>
  )
}

export const getStaticPaths: GetStaticPaths<{ id: string }> = async () => {
  const posts = await getPostsLocal()
  return {
    paths: posts.map((post) => {
      return { params: { id: post.id.toString() } }
    }),
    // TODO: maybe remove
    fallback: "blocking", // fallback tries to regenerate ArtistPage if Artist did not exist during building
  }
}

export const getStaticProps: GetStaticProps<{ post: Post }, { id: string }> = async ({
  params,
}) => {
  if (params == undefined) {
    return { notFound: true }
  }
  const post = await getPostLocal(+params.id)
  return post ? { props: { post: post } } : { notFound: true }
}
export default PostPage
