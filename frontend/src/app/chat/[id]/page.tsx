import { redirect } from "next/navigation"

import ChatPage, { ChatPageFirstMessage } from "@/app/chat/[id]/chat-page"
import { getUserLocal } from "@/api/users/edge"

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

const Page = ({ params }: Props) => {
  if (params.id.startsWith("u")) {
    const userId = +params.id.slice(1)
    return <ChatPageFirstMessage userId={userId} />
  }

  const id = +params.id
  if (isNaN(id) || id < 0) {
    return redirect("/chat")
  }

  return <ChatPage id={id} />
}

export default Page
