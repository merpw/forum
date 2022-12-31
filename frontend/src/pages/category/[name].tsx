import { GetStaticPaths, GetStaticProps, NextPage } from "next"
import Head from "next/head"
import { getCategoriesLocal, getCategoryPostsLocal } from "../../api/posts/fetch"
import { PostList } from "../../components/posts/list"
import { Post } from "../../custom"

const CategoryPage: NextPage<{ category_name: string; posts: Post[] }> = ({
  category_name,
  posts,
}) => {
  return (
    <>
      <Head>
        <title>{`${category_name} - Forum`}</title>
        <meta name={"description"} content={"The friendliest forum"} />
        <meta property={"og:title"} content={"FORUM"} key={"title"} />
        <meta name={"og:description"} content={"The friendliest forum"} />
      </Head>
      <h1 className={"text-2xl font-light mb-5"}>
        Category <span className={"font-normal"}>{category_name}</span>
      </h1>
      <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
    </>
  )
}

export const getStaticPaths: GetStaticPaths<{ name: string }> = async () => {
  const categories = await getCategoriesLocal()
  return {
    paths: categories.map((category) => {
      return { params: { name: category } }
    }),
    // TODO: maybe remove
    fallback: "blocking", // fallback tries to regenerate ArtistPage if Artist did not exist during building
  }
}

export const getStaticProps: GetStaticProps<{ posts: Post[] }, { name: string }> = async ({
  params,
}) => {
  if (params == undefined) return { notFound: true }

  let category_name = params.name.toLowerCase()

  const { posts } = await getCategoryPostsLocal(category_name)
  if (posts == undefined) return { notFound: true }

  // capitalize
  category_name = category_name.charAt(0).toUpperCase() + category_name.slice(1)

  return {
    props: { posts: posts, category_name: category_name },
    revalidate: 10,
  }
}

export default CategoryPage
