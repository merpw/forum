import { Message } from "../../types"

import { List } from "./list.js"
import { Chat } from "./chat.js"
import { renderChatMessages } from "./helpers/events"

export class Client {
  activeChat: Chat | null 
  chatIds: number[] | undefined
  chats: Map<number, Chat> // chatId => Chat
  messages: Map<number, Message | null> // messageId => Message
  userChats: Map<number, number | null> // userId => chatId
  usersOnline: number[] | undefined


  constructor(){
    this.activeChat = null 
    this.chatIds = undefined
    this.chats = new Map<number, Chat>
    this.messages = new Map<number, Message | null>
    this.userChats = new Map<number, number | null>
    this.usersOnline = undefined
    window.addEventListener("renderChatList", () => {
      setTimeout(() => {
        this.list = new List(this.onlineUsers, this.sortedChats, this.chatUserIds)
        this.list.render()
      }, 100)
    }) 
  }

  list: List | undefined

  reset(){
    this.activeChat = null
    this.chatIds = undefined
    this.chats.clear()
    this.messages.clear()
    this.userChats.clear()
    this.usersOnline = undefined
  }

  private get sortedChats(): Chat[]{
    return [...this.chats.values()].sort((a, b) => a.lastMessageId - b.lastMessageId).reverse()
  }

  get chatUserIds(): number[]{
    return this.sortedChats.map((chat) => chat.userId)
  }

  get onlineUsers(): number[]{
    return this.usersOnline as number[] 
  }

  set onlineUsers(users: number[]){
    this.usersOnline = users 
  }
}


