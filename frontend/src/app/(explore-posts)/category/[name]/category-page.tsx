"use client"

import { NextPage } from "next"

import { PostList } from "@/components/posts/list"
import { useCategoryPosts } from "@/api/posts/hooks"

const CategoryPage: NextPage<{ category: string }> = ({ category }) => {
  const { posts } = useCategoryPosts(category)

  if (!posts) {
    return null
  }

  return (
    <>
      <PostList posts={posts.sort((a, b) => b.date.localeCompare(a.date))} />
    </>
  )
}

export default CategoryPage
