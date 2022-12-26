import {Post} from "../../custom";
import {posts} from "../dummy";
import useSWR from "swr";

export const getPostsLocal = (): Promise<Post[]> => Promise.resolve(posts)

export const getPostLocal = (id: number): Promise<Post | undefined> =>
    Promise.resolve(posts.find((post) => post.id == id))

export const getUserPosts = (user_id: number): Promise<{ posts: Post[] }> => {
    return Promise.resolve({ posts: posts.filter((post) => post.author.id == user_id) })
}

export const useUserPosts = (user_id: number | undefined) => {
    const { data, error } = useSWR<{ posts: Post[] }>(
        user_id ? `/api/user/${user_id}/posts` : null,
        () => getUserPosts(user_id || -1)
    )

    return {
        isError: error != undefined,
        isLoading: !error && !data,
        posts: data?.posts,
    }
}