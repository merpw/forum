import { client } from "../../main.js"
import { User } from "../../types.js"
import { state } from "../authorization/auth.js"
import { Chat } from "./chat.js"
import { sendWsObject } from "./helpers/sendobject.js"
import { createElement } from "../utils.js"

export class List {
  public users: User[]
  public chats: Chat[]
  private onlineList: number[]

  constructor(onlineList: number[], chats: Chat[], chatUserIds: number[]) {
    const users = [...state.users.values()]
      .filter((user) => !chatUserIds.includes(user.id))

    chats.forEach((chat) => {
      chat.online = onlineList.includes(chat.userId) ? true : false
    })

    this.onlineList = onlineList
    this.users = users
    this.chats = chats

    this.render()
  }

  chatsList = document.getElementById("your-chats-list") as HTMLUListElement
  usersList = document.getElementById(
    "online-users-list"
  ) as HTMLUListElement

  private render() {
    // Add chats to DOM.
    this.chatsList.replaceChildren()
    for (const chat of this.chats) {
      const chatUserElement = createElement(
        "li",
        "chat-user",
        `chatId${chat.chatId}`
      )
      const userNameElement = createElement("p", null, null, chat.name)
      chatUserElement.appendChild(userNameElement)

      const notification = createElement("i", "bx bx-message-dots")
      const amount = createElement("p", null, null, `${chat.unreadMessages}`)
      if (chat.unreadMessages === 0) {
        notification.style.display = "none"
        amount.style.display = "none"
      } else {
        notification.style.display = "flex"
        amount.style.display = "flex"
      }

      chatUserElement.append(notification, amount)

      if (chat.online) {
        chatUserElement.classList.add("online")
      } else {
        chatUserElement.classList.add("offline")
      }

      chatUserElement.addEventListener("click", (e) => {
        if (chat.chatId === client.activeChat?.chatId) {
          return
        }

        chat.open()
        e.preventDefault()
      })

      this.chatsList.appendChild(chatUserElement)
    }

    // Add online users to DOM.
    this.usersList.replaceChildren()

    for (const user of this.users) {
      const userElement = createElement("li", "chat-user")
      const nameElement = createElement("p", null, null, user.name)
      userElement.appendChild(nameElement)

      userElement.addEventListener("click", (e) => {
        this.createChat(user.id)
        setTimeout(() => {
          const chatId = client.userChats.get(user.id) as number
          const chat = client.chats.get(chatId) as Chat
          chat.open()
        }, 100)
        e.preventDefault()
      })

      if(this.onlineList.includes(user.id)){
        userElement.classList.add("online")
      } else {
        userElement.classList.add("offline")
      }

      this.usersList.appendChild(userElement)
    }
  }

  private createChat(userId: number) {
    sendWsObject({
      type: "post",
      item: {
        url: "/chat/create",
        data: {
          userId: userId,
        },
      },
    })
  }
}
