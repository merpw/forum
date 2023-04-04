import { GetStaticPaths, GetStaticProps, NextPage } from "next"
import { Comment, Post } from "@/custom"

import Link from "next/link"
import { FC, useEffect, useState } from "react"
import Head from "next/head"

import { getPostCommentsLocal, getPostLocal, getPostsLocal } from "@/api/posts/fetch"
import { useMe } from "@/api/auth"
import { CreateComment, useComments } from "@/api/posts/comment"
import { FormError } from "@/components/error"
import { Category, ReactionsButtons, ReactionsCommentButtons } from "@/components/posts/reactions"
import { SWRConfig, SWRConfiguration, unstable_serialize } from "swr"
import ReactTextareaAutosize from "react-textarea-autosize"
import useDates from "@/helpers/dates"

const PostPage: NextPage<{ post: Post; fallback: SWRConfiguration }> = ({ post, fallback }) => {
  const { localDate, relativeDate } = useDates(post.date)

  return (
    <SWRConfig value={fallback}>
      <Head>
        <title>{`${post.title} - Forum`}</title>
        <meta property={"og:title"} content={`${post.title} - Forum`} />

        {/* TODO: change to description */}
        <meta name={"description"} content={post.content.slice(0, 200)} />
        <meta property={"og:description"} content={post.content.slice(0, 200)} />
      </Head>
      <div className={"m-5"}>
        <div className={"mb-3"}>
          <h1 className={"text-3xl mb-2 "}>{post.title}</h1>
          <hr />
        </div>
        <p className={"whitespace-pre-line"}>{post.content}</p>
        <hr className={"mt-4"} />
        <div className={"border-t py-2 flex flex-wrap"}>
          <ReactionsButtons post={post} />
          <Category post={post} />

          <span className={"ml-auto"}>
            <span title={localDate}>{relativeDate}</span>
            {" by "}
            <span className={"text-xl hover:opacity-50"}>
              <Link href={`/user/${post.author.id}`}>{post.author.name}</Link>
            </span>
          </span>
        </div>
        <div>
          <h2 className={"text-2xl my-4"}>Comments:</h2>
          <CommentForm post={post} />
          <Comments post={post} />
        </div>
      </div>
    </SWRConfig>
  )
}

const CommentForm: FC<{ post: Post }> = ({ post }) => {
  const { isLoggedIn } = useMe()
  const { mutate: mutateComments } = useComments(post.id)

  const [text, setText] = useState("")
  const [formError, setFormError] = useState<string | null>(null)

  const [isSame, setIsSame] = useState(false)

  useEffect(() => setIsSame(false), [text])

  if (!isLoggedIn) return null

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()

        if (isSame) return

        setIsSame(true)
        if (formError != null) setFormError(null)

        CreateComment(post.id, text)
          .then(() => {
            setText("")
            mutateComments()
          })
          .catch((err) => {
            if (err.code == "ERR_BAD_REQUEST") {
              setFormError(err.response?.data as string)
            } else {
              // TODO: unexpected error
            }
          })
      }}
      className={"mb-5"}
    >
      <div className={"mb-3"}>
        <label
          htmlFor={"comment-text"}
          className={"block mb-2 text-sm font-medium text-gray-900 dark:text-white"}
        >
          Write a comment:
        </label>

        <ReactTextareaAutosize
          id={"comment-text"}
          className={
            "w-full bg-gray-50 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 "
          }
          value={text}
          onInput={(e) => setText(e.currentTarget.value)}
          required
        />
      </div>

      <FormError error={formError} />

      <button
        type={"submit"}
        className={
          "text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
        }
      >
        Submit
      </button>
    </form>
  )
}

const Comments: FC<{ post: Post }> = ({ post }) => {
  const { comments } = useComments(post.id)

  return (
    <div className={"flex flex-col gap-3"}>
      {comments
        ?.sort((a, b) => b.date.localeCompare(a.date))
        .map((comment, key) => (
          <CommentCard comment={comment} post={post} key={key} />
        ))}
    </div>
  )
}

const CommentCard: FC<{ comment: Comment; post: Post }> = ({ comment, post }) => {
  const { localDate, relativeDate } = useDates(comment.date)
  return (
    <div className={"border rounded p-5"}>
      <Link href={`/user/${comment.author.id}`}>
        <h3 className={"text-lg hover:opacity-50"}>{comment.author.name}</h3>
      </Link>
      <p className={"whitespace-pre-line"}>{comment.content}</p>
      <hr className={"mt-4 mb-2"}></hr>
      <span className={"flex"}>
        <ReactionsCommentButtons post={post} comment={comment} />

        <span suppressHydrationWarning className={"ml-auto"} title={localDate}>
          {relativeDate}
        </span>
      </span>
    </div>
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
  const comments = await getPostCommentsLocal(+params.id)

  return post
    ? {
        props: {
          post: post,
          fallback: {
            [unstable_serialize(["api", "posts", post.id, "comments"])]: comments,
          },
        },
        revalidate: 1,
      }
    : { notFound: true, revalidate: 1 }
}
export default PostPage
