import { ChatUser, Message } from "../../types"

import { getUserById } from "../../api/get.js"
import { client } from "../../main.js"
import { createElement } from "../utils.js"

import { sendWsObject } from "./helpers/sendobject.js"

export class Client {
  activeChat: Chat | null // signals which 
  chatIds: number[] | undefined
  chats: Map<number, Chat | null> // chatId => Chat
  messages: Map<number, Message | null> // messageId => Message
  userChats: Map<number, number | null> // userId => chatId
  usersOnline: number[] | undefined


  constructor(){
    this.activeChat = null 
    this.chatIds = undefined
    this.chats = new Map<number, Chat | null>
    this.messages = new Map<number, Message | null>
    this.userChats = new Map<number, number | null>
    this.usersOnline = undefined
  }

  reset(){
    this.activeChat = null
    this.chatIds = undefined
    this.chats.clear()
    this.messages.clear()
    this.userChats.clear()
    this.usersOnline = undefined
  }

  getAllChats(){
    sendWsObject({
      type: "get",
      item: {
        url: "/chat/all"
      }

    })

    window.addEventListener('renderChatList', () => {
      console.log(client)
      // renderChatList()
    })
  }
}


// Manage state in the chat list with ez pz method managing xD
export class List {
  // online: ChatUser[] = []
  // chats: Chat[] = []

  constructor(
    public online: ChatUser[],
    public chats: Chat[]
  ){}

  async getOnlineUsers(){
    const onlineList: ChatUser[] = []

    for (const id of client.usersOnline as number[]){
      
      if(client.userChats.has(id)){
        const chatId = client.userChats.get(id) as number
        this.chats.push(client.chats.get(chatId) as Chat)
      } else {
        onlineList.push(await getUserById(id))
      }

    }

    this.online = onlineList
  }

  sortChats(){
    this.chats.sort((a, b) =>  a.lastMessageId - b.lastMessageId)
  }
}

export class Chat {
  constructor(
    readonly name: string, 
    readonly userId: number, 
    readonly chatId: number, 
    public lastMessageId: number
  ){this.getMessageIds()}

  messageIds: number[] = [] // Set with an easy call to WS
  messages: Message[] = [] // Unshifted to when open chat
  
  private range = 10 // Determines how many messages should be loaded



  private getMessageIds(){
    sendWsObject({
      type: "get",
      item: {
        url: `/chat/${this.chatId}/messages`
      }
    })
  }

  open(){
    client.activeChat = this // SET activeChat to THIS chat
    this.getMessageIds()

    // const chatArea = document.getElementById("chat-area") as HTMLDivElement
    // const chatWindow = createElement("div", "chat show-chat")

    const chatHeader = createElement(
      "div",
      "chat-header"
    )

    const chatName = createElement(
      "div",
      "chat-name",
      null,
      this.name
    ) as HTMLDivElement

    const closeBtn = createElement("a", "closebtn", "chat-window-close", null, "<i class='bx bx-x'>")

    closeBtn.addEventListener("click", (e) => {
      e.preventDefault()
      this.close()
    })
    chatHeader.append(chatName, closeBtn)

    const chatMessages = createElement("div", "chat-messages", `Chat${this.chatId}`)

    // const chatFormContainer = createElement("div", "chat-form-container")
    const chatForm = createElement("form", "chat-form")
    chatForm.setAttribute("autocomplete", "off")
    const messageField = createElement("input", null, "chat-text")
    messageField.setAttribute("maxlength", "150")

    chatMessages.addEventListener("messageEvent", () => {
      // updateMessagges(this.chatId)
    })

    const sendButton = document.getElementById("chat-send") as HTMLElement
    sendButton.addEventListener('sendMessage', (e) => {
      const content = document.getElementById("chat-content") as HTMLInputElement
      this.sendMessage(content.value)
      e.preventDefault()
    })

  }

  close(){
    client.activeChat = null
  }

  private sendMessage(content: string){
    // TODO: Add logic for sending messages HERE.
    sendWsObject({
      type: "post",
      item: {
        url: `/chat/${this.chatId}/message`,
        data: {
          content: content
        }

      }
    })
  }

  // private updateMessages(chatId: number){

  // }

}

// function createChat(userId: number){
//   sendWsObject({
//     type: "post",
//     item: {
//       url: "/chat/create",
//       data: {
//         userId: userId
//       }
//     }
//   })
// }
