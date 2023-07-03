import { ChatUser, Message } from "../types"
import { sendWsObject } from "./ws.js"
import { userInfo } from "./auth.js"
import { createElement, iterator } from "./utils.js"
import { getUserById, getUserIds } from "../api/get.js"

export const chat = {
  list: {
    idMap: new Map<number, number>, // companionId, chatId
    users: [] as ChatUser[],
  }, 

  messages: {
    list: [] as Message[],
    buffer: [] as Message[],
    current: {} as Message
  },

  current: {
    username: "",
    userId: -1,
    chatId: -1,
  }
}

chat.current

export const chatList = {
  Ids: new Map<number, number>, //companionId, chatId
  Users: [] as ChatUser[],
}

export const messages = {
  list: [] as Message[],
  buffer: [] as Message[],
  current: {} as Message 
}


export const currentChat = {
  username: "",
  userId: -1,
  chatId: -1
}

const getChatUsers = async () => {
  Object.assign(chatList, {
    Ids: new Map<number, number>,
    Users: [],
  })

  const userIds: iterator = await getUserIds()
  for (const id of Object.values(userIds)) {
    if (id !== userInfo.Id) {
      chatList.Users.push(await getUserById(id))
    }
  }

  sendWsObject({
    type: "get",
    item: {
      url: "/users/online"
    }
  })

  for (const user of chatList.Users) {
    if (!chatList.Ids.has(user.Id)) {
      sendWsObject({
        type: "post",
        item: {
          url: "/chat/create",
          data: {
            userId: user.Id,
          },
        },
      }) 
    }
  }
}

// Display the current chat based on chatId
const showChat = async (id: number) => {
  const chatArea = document.getElementById("chat-area") as HTMLDivElement
  const chat = createElement("div", "chat show-chat")

  for (const user of chatList.Users) {
    if (user.Id === id && user.UnreadMsg) {
      user.UnreadMsg = false
      renderChatList()
      break
    }
  }

  const chatId = chatList.Ids.get(id) as number
  getMessages(chatId)
  const user = await getUserById(id.toString())

  Object.assign(currentChat, {
    username: user.Name,
    userId: user.Id,
    chatId: chatId,
  })

  const chatName = createElement(
    "div",
    "chat-name",
    null,
    user.Name
  ) as HTMLDivElement

  chatName.style.paddingBottom = "10px"

  const chatMessages = createElement("div", "chat-messages", `Chat${chatId}`)

  const chatFormContainer = createElement("div", "chat-form-container")
  const chatForm = createElement("form", "chat-form")
  const messageField = createElement("input", null, "chat-text")
  messageField.setAttribute("maxlength", "150")
  chatMessages.addEventListener("messageEvent", () => {
    updateChat(chatId)
  })

  const messageSend = createElement("button", null, "chat-send", "Send")
  messageSend.style.marginLeft = "4px"
  messageSend.addEventListener("click", (e) => {
    sendMessage(chatId)
    e.preventDefault()
  })

  chatForm.append(messageField, messageSend)
  chatFormContainer.appendChild(chatForm)

  setTimeout(() => {
    for (const message of Object.values(messages.list)) {
      const msgElement = createElement("div", "message") as HTMLDivElement
      const date = new Date(message.timestamp)
      const formatDate = date
      .toLocaleString("en-GB", { timeZone: "EET" })
      .slice(0, 10)
      if (message.authorId === userInfo.Id) {
        msgElement.classList.add("send")
        msgElement.textContent = "You:\n" + message.content + "\n" + formatDate
      } else if (message.authorId === -1) {
        msgElement.classList.add("status")
        msgElement.textContent = message.content + "\n" + formatDate
      } else {
        msgElement.classList.add("recieve")
        msgElement.textContent = currentChat.username + ":\n" + message.content + "\n" + formatDate
      }
      chatMessages.appendChild(msgElement)
    }

    chat.append(chatName, chatMessages, chatFormContainer)
    chatArea.replaceChildren(chat)
  }, 200)

}

const sendMessage = async (chatId: number): Promise<void> => {
  return new Promise<void>((resolve) => {
  const content = document.getElementById("chat-text") as HTMLInputElement
  let message = content.value.toString().trim()
  if (message.length > 0) {
    content.value = ""
    if (message.length > 150) {
      message = message.slice(0, 150)
    }
    sendWsObject({
      type: "post",
      item: {
        url: `/chat/${chatId}/message`,
        data: {
          content: message,
        },
      },
    })
  }
    setTimeout(resolve, 10)
  })
}

export const updateChat = (chatId: number): Promise<void> => {
  return new Promise<void>((resolve) => {
    const message = messages.list[0]
    setTimeout(() => {
      const chatMessages = document.getElementById(
        `Chat${chatId}`
      ) as HTMLDivElement
      const msgElement = createElement(
        "div",
        "message",
        null,
        message.content
      ) as HTMLDivElement
      const date = new Date(message.timestamp)
      const formatDate = date
      .toLocaleString("en-GB", { timeZone: "EET" })
      .slice(0, 10)
      if (message.authorId === userInfo.Id) {
        msgElement.classList.add("send")
        msgElement.textContent = "You:\n" + message.content + "\n" + formatDate
      } else if (message.authorId === -1) {
        msgElement.classList.add("status")
        msgElement.textContent = message.content + "\n" + formatDate
      } else {
        msgElement.textContent = currentChat.username + ":\n" + message.content + "\n" + formatDate
        msgElement.classList.add("recieve")
      }
      if (currentChat.chatId === messages.current.chatId) {
        chatMessages.prepend(msgElement)
      } 
    }, 5)
    resolve(renderChatList())
  })
}

const getMessages = async (chatId: number): Promise<void> => {
  return new Promise<void>((resolve, reject) => {
    // Reset the state of the messages
    Object.assign(messages, { list: [] })
    setTimeout(() => {
      sendWsObject({
        type: "get",
        item: {
          url: `/chat/${chatId}/messages`,
        },
      })
    }, 20)
    setTimeout(() => {
      if (messages.list.length === 0) {
        reject("Empty chat")
      } else {
        resolve()
      }
    }, 50)
  })
}

export const displayChatUsers = async () => {
  const onlineTitle = document.getElementById("online-title") as HTMLElement
  const offlineTitle = document.getElementById("offline-title") as HTMLElement

  /* Add eventlisteners to show/hide users in chat list */
  onlineTitle.addEventListener("click", toggleOnline)
  offlineTitle.addEventListener("click", toggleOffline)

  await getChatUsers()
  renderChatList()
}

export function renderChatList(): Promise<void> {
  return new Promise((resolve) => {
    sendWsObject({
      type: "get",
      item: {
        url: "/chat/all",
      },
    }) 
    setTimeout(() => {
      const onlineList = document.getElementById("online-users") as HTMLUListElement,
      offlineList = document.getElementById("offline-users") as HTMLUListElement

      /* Resets the list if user logging in or out */
      onlineList.replaceChildren()
      offlineList.replaceChildren()

      for (const u of chatList.Users) {
        const user = createElement("li", "chat-user", u.Name)
        const name = createElement("p", null, null, u.Name)
        user.appendChild(name)

        if (u.UnreadMsg) {
          const unread = createElement("i", "bx bx-message-dots")
          user.appendChild(unread)
        }

        user.addEventListener("click", () => {
          showChat(u.Id)
        })

        if (u.Online) {
          user.classList.add("online")
          onlineList.appendChild(user)
        } else {
          user.classList.add("offline")
          offlineList.appendChild(user)
        }
      }
    }, 120)
    setTimeout(resolve, 130)
  })
}


const toggleOnline = () => {
  const onlineToggle = document.getElementById("online-toggle") as HTMLElement,
  onlineUsers = document.getElementById("online-users") as HTMLUListElement
  if (onlineToggle.className === "bx bx-chevron-down") {
    onlineToggle.className = "bx bx-chevron-right"
    onlineUsers.style.display = "none"
  } else {
    onlineToggle.className = "bx bx-chevron-down"
    onlineUsers.style.display = "block"
  }
}

const toggleOffline = () => {
  const offlineToggle = document.getElementById("offline-toggle") as HTMLElement
  const offlineUsers = document.getElementById(
    "offline-users"
  ) as HTMLUListElement
  if (!offlineToggle) return

  if (offlineToggle.className === "bx bx-chevron-down") {
    offlineToggle.className = "bx bx-chevron-right"
    offlineUsers.style.display = "none"
  } else {
    offlineToggle.className = "bx bx-chevron-down"
    offlineUsers.style.display = "block"
  }
}
