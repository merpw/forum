import { Message } from "../../types"

import { sendWsObject } from "./helpers/sendobject.js"

export class Client {
  activeChatId: number | null
  chatIds: number[] | undefined
  chats: Map<number, Chat | null>
  chatMessages: Map<number, number[]>
  messages: Map<number, Message | null>
  userChats: Map<number, number | null>
  usersOnline: number[] | undefined  


  constructor(){
    this.activeChatId = null 
    this.chatIds = undefined
    this.chats = new Map<number, Chat | null>
    this.chatMessages = new Map<number, number[]>
    this.messages = new Map<number, Message | null>
    this.userChats = new Map<number, number | null>
    this.usersOnline = undefined
  }
  
  private initialize(){
    sendWsObject({
      type: "get",
      item: {
        url: "/chat/all"
      }

    })

    window.addEventListener('renderChatList', () => {
      // renderChatList()
    })

    window.addEventListener("createChat", () => {
      // createChat()
    })


  }
}

export class Chat {
  userId: number
  chatId: number
  lastMessageId: number
  messageIds: number[]
  messages: Map<number, Message>

  constructor(userId: number, chatId: number, lastMessageId: number)  {
    this.userId = userId
    this.chatId = chatId
    this.lastMessageId = lastMessageId
    this.messageIds = [lastMessageId]
    this.messages = new Map<number, Message>

    window.addEventListener('renderChatWindow', (e) => {
      // renderChat(chatId)
      e.preventDefault()
    })
    
    window.addEventListener('sendMessage', (e) => {
      // sendMessage(chatId)
      e.preventDefault()
    })
  }
}


