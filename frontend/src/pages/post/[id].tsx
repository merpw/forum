import { GetStaticPaths, GetStaticProps, NextPage } from "next"
import { Post } from "../../custom"

import Link from "next/link"
import { FC, useEffect, useState } from "react"
import Head from "next/head"
import moment from "moment"
import { getPostLocal, getPostsLocal } from "../../api/posts/fetch"
import { useMe } from "../../api/auth"
import { CreateComment } from "../../api/posts/comment"
import { ReactionsButtons } from "../../components/posts/reactions"
import { FormError } from "../../components/error"

const PostPage: NextPage<{ post: Post }> = ({ post }) => {
  return (
    <>
      <Head>
        <title>{`${post.title} - Forum`}</title>
      </Head>
      <div className={"m-5"}>
        <div className={"mb-3"}>
          <h1 className={"text-3xl mb-2 "}>{post.title}</h1>
          <hr />
        </div>
        <p>{post.content}</p>
        <hr className={"mt-4"} />
        <div className={"border-t py-2 flex flex-wrap"}>
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
        <div>
          <h2 className={"text-2xl my-4"}>Comments:</h2>
          <CommentForm post={post} />
          <Comments post={post} />
        </div>
      </div>
    </>
  )
}

const CommentForm: FC<{ post: Post }> = ({ post }) => {
  const { isLoggedIn } = useMe()
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

        if (formError != null) setFormError(null)

        CreateComment(post.id, text)
          .then((id) => {
            console.log(id)
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

        <textarea
          id={"comment-text"}
          className={
            "w-full bg-gray-50 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 "
          }
          rows={text.split("\n").length}
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
  if (post.comments.length == 0) {
    return <div>There are no comments yet, write one first!</div>
  }

  return (
    <div className={"flex flex-col gap-3"}>
      {post.comments.map((comment, key) => (
        <div className={"border rounded p-5"} key={key}>
          <Link href={`/user/${comment.author.id}`}>
            <h3 className={"text-lg hover:opacity-50"}>{comment.author.name}</h3>
          </Link>
          <p>{comment.text}</p>
          <hr className={"mt-4 mb-2"}></hr>
          <span className={"flex"}>
            <ReactionsButtons post={post} comment={comment} />

            <span
              className={"ml-auto"}
              title={moment(comment.date).local().format("DD.MM.YYYY HH:mm:ss")}
            >
              {moment(comment.date).fromNow()}
            </span>
          </span>
        </div>
      ))}
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
  return post ? { props: { post: post }, revalidate: 10 } : { notFound: true }
}
export default PostPage