import { useEffect } from "react"

import { useAppDispatch, useAppSelector } from "@/store/hooks"
import wsActions from "@/store/ws/actions"

export const useChatIds = () => {
  const chatIds = useAppSelector((state) => state.chats.chatIds)
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chatIds) {
      dispatch(wsActions.sendGet("/chat/all"))
    }
  }, [chatIds, dispatch])

  return { chatIds }
}

export const useChat = (id: number) => {
  const chat = useAppSelector((state) => state.chats.chats?.[id])
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chat) {
      dispatch(wsActions.sendGet(`/chat/${id}`))
    }
  }, [chat, dispatch, id])

  return { chat }
}

type AssociationData = { companionId: number } | { groupId: number }

export const useCreateChat = () => {
  const dispatch = useAppDispatch()

  return (association: AssociationData) => {
    dispatch(wsActions.sendPost("/chat/create", association))
  }
}

export const useAssociatedChat = (association: AssociationData) => {
  const isUser = "companionId" in association
  const chatId = useAppSelector((state) =>
    isUser
      ? state.chats.userChats?.[association.companionId]
      : state.chats.groupChats?.[association.groupId]
  )
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chatId) {
      dispatch(
        wsActions.sendGet(
          isUser ? `/users/${association.companionId}/chat` : `/groups/${association.groupId}/chat`
        )
      )
    }
  }, [association, chatId, dispatch, isUser])

  return { chatId }
}

export const useGroupChat = (groupId: number) => {
  const chatId = useAppSelector((state) => state.chats.groupChats?.[groupId])
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (!chatId) {
      dispatch(wsActions.sendGet(`/groups/${groupId}/chat`))
    }
  }, [chatId, dispatch, groupId])

  return { chatId }
}
