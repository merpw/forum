import { GetStaticProps, NextPage } from "next"
import { getPostsLocal } from "../fetch/server-side"

import Head from "next/head"
import { Post } from "../custom"
import Link from "next/link"

const Home: NextPage<{ posts: Post[] }> = ({ posts }) => {
  return (
    <>
      <Head>
        <title>FORUM</title>
        <meta name={"description"} content={"The friendliest forum"} />
        <meta property={"og:title"} content={"FORUM"} key={"title"} />
        <meta name={"og:description"} content={"The friendliest forum"} />
      </Head>
      <div>
        {posts.map((post, key) => {
          return PostCard(post, key)
        })}
      </div>
    </>
  )
}

const PostCard = (post: Post, key: number) => (
  <div key={key} className={"border rounded-lg sm:m-5 p-7"}>
    <div className={"mb-3"}>
      <Link href={`/post/${post.id}`}>
        <h1 className={"text-2xl"}>{post.title}</h1>
      </Link>
      <hr className={"mt-2"} />
    </div>
    <p>{post.content}</p>
    <hr className={"mt-4"} />
    <div>
      <span>{new Date(Date.parse(post.date)).toLocaleString("fi")}</span> by{" "}
      <Link href={`/user/${post.author.id}`}>
        <span className={"text-xl"}>{post.author.name}</span>
      </Link>
    </div>
  </div>
  // <div className={"mx-auto my-2 text-center p-1 w-max"} key={key}>
  //   <Link href={`/artist/${post.name}`}>
  //     <Image
  //       src={process.env.NEXT_PUBLIC_GROUPIE_BACKEND_HOST + post.image}
  //       alt={`Image of ${post.name}`}
  //       width={240}
  //       height={240}
  //       blurDataURL={"placeholder"} // TODO: add base64 placeholder
  //       placeholder={"blur"}
  //       className={"rounded-full hover:brightness-125"}
  //     />
  //     <p className={"text-xl"}>{post.name}</p>
  //   </Link>
  // </div>
)

export const getStaticProps: GetStaticProps<{ posts: Post[] }> = async () => {
  const posts: Post[] = await getPostsLocal()

  return {
    props: { posts: posts },
    revalidate: 60,
  }
}

export default Home
