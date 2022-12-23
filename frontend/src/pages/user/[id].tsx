import { GetStaticPaths, GetStaticProps, NextPage } from "next"
import { User } from "../../custom"

import { getUserLocal, getUsersLocal } from "../../fetch/server-side"

const UserPage: NextPage<{ user: User }> = ({ user }) => {
  return (
    <div>
      <h1 className={"text-xl"}>{user.name}</h1>
    </div>
  )
}

export const getStaticPaths: GetStaticPaths<{ id: string }> = async () => {
  const users = await getUsersLocal()
  return {
    paths: users.map((user) => {
      return { params: { id: user.id.toString() } }
    }),
    // TODO: maybe remove
    fallback: "blocking", // fallback tries to regenerate ArtistPage if Artist did not exist during building
  }
}

export const getStaticProps: GetStaticProps<{ user: User }, { id: string }> = async ({
  params,
}) => {
  if (params == undefined) {
    return { notFound: true }
  }
  const user = await getUserLocal(+params.id)
  return user ? { props: { user: user } } : { notFound: true }
}
export default UserPage
