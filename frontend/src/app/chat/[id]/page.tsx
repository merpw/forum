import { redirect } from "next/navigation"
import { cookies } from "next/headers"

import ChatPage, { ChatPageFirstMessage } from "@/app/chat/[id]/chat-page"
import { getUserLocal } from "@/api/users/edge"
import checkSession from "@/api/auth/edge"

type Props = { params: { id: string } }

export const generateMetadata = async ({ params }: Props) => {
  if (params.id.startsWith("u")) {
    const userId = +params.id.slice(1)
    if (isNaN(userId) || userId < 0) {
      return redirect("/chat")
    }

    try {
      const user = await getUserLocal(userId)
      return {
        title: `Chat with ${user.username}`,
      }
    } catch (error) {
      return redirect("/chat")
    }
  }

  return {
    title: `Chats`,
  }
}

const Page = async ({ params }: Props) => {
  if (params.id.startsWith("u")) {
    const userId = +params.id.slice(1)
    const token = cookies().get("forum-token")?.value
    if (!token) {
      return redirect("/login")
    }
    try {
      const requestUserId = await checkSession(token)
      if (requestUserId === userId) {
        return redirect("/chat")
      }
      return <ChatPageFirstMessage userId={userId} />
    } catch (error) {
      return redirect("/chat")
    }
  }

  const id = +params.id
  if (isNaN(id) || id < 0) {
    return redirect("/chat")
  }

  return <ChatPage id={id} />
}

export default Page
