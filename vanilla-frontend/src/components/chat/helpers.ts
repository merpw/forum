import { getUserIds, getUserById } from "../../api/get.js"

import { userInfo } from "../authorization/auth.js"

import { chatList, renderChatList, messages, currentChat } from "./chat.js"
import { messageEvent, ws } from "./ws.js"

import { createElement, iterator } from "../utils.js"
import { Message } from "../../types.js"


/* WS FUNCTIONS */

// Updates the state of online users in the chat list
export async function updateOnlineUsers(users: number[]): Promise<void> {
    sendWsObject({
      type: "get",
      item: {
        url: "/chat/all",
      },
    })

    Object.assign(chatList, {
      Ids: new Map<number, number>(),
      Users: [],
    })
  
    const userIds: iterator = await getUserIds()

    for (const id of Object.values(userIds)) {
      if (id !== userInfo.Id) {
        chatList.Users.push(await getUserById(id))
      }
    }
  
    for (const user of Object.values(chatList.Users)) {
      user.Online = users.includes(user.Id)
    }
  
    renderChatList()
  }
  

  
  // Loops through the ID array of all chats and sends a message to get matching chatId/userId pairs
  export const getChatIds = (resp: iterator): void => {
    for (const chat of resp.item.data) {
      if (!chatList.Ids.has(chat)) {
        sendWsObject({
          type: "get",
          item: {
            url: `/chat/${chat}`,
          },
        })
      }
    }
    return
  }
  
  // Gets one message with all its information
  export const getMessage = (id: number): void => {
    sendWsObject({
      type: "get",
      item: {
        url: `/message/${id}`,
      },
    })
    setTimeout(() => {
      const chat = document.getElementById(`Chat${currentChat.chatId}`) as HTMLDivElement
      if (chat) {
        // If the message is in the opened chat, dispatch the messageEvent that updates the chat with newest message.
        chat.dispatchEvent(messageEvent)
      } else {
        // Otherwise, update unread messages in the chat list
        updateUnreadMessages()
      }
    }, 80)
  }
  
  // Gets a list of all messages in the provided ID array. Currently getting ALL messages.
  export const getMessageList = (ids: number[]): void => {
    for (const id of ids) {
        sendWsObject({
        type: "get",
        item: {
          url: `/message/${id}`,
        },
      })
    }
  }

  export const sendMessage = (chatId: number): void => {
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
    renderChatList()
  }
  
  // Sends a WebSocket object as a message to the WebSocket server
  export function sendWsObject(obj: object): void {
      if (ws.readyState === 1) {
        ws.send(JSON.stringify(obj))
      }
  }



/* CHAT STATE/DOM FUNCTIONS */
  
// Updates the state of unread messages in the chat list
  export function updateUnreadMessages(): void {
    loop:
    for (const [userId, chatId] of chatList.Ids) {
      if (chatId === messages.current.chatId) {
        for (const user of Object.values(chatList.Users)) {
          if (user.Id === userId) {
            user.UnreadMsg = true
            break loop
          }
        }
        break
      }
    }

    renderChatList()
  }

  // Creates a message in the DOM and returns the HTMLDivElement
  export const createMessage = (message: Message): HTMLDivElement => {
    const msgElement = createElement("div", "message") as HTMLDivElement
    const date = new Date(message.timestamp)
    const formatDate = date
    .toLocaleString("en-GB", { timeZone: "EET" })
    .slice(0, 10)
    const dateElement = createElement("p", "date", null, formatDate)
    if (message.authorId === userInfo.Id) {
      msgElement.classList.add("send")
      msgElement.textContent = "You:\n" + message.content + "\n"
    } else if (message.authorId === -1) {
      msgElement.classList.add("status")
      msgElement.textContent = message.content + "\n"
    } else {
      msgElement.classList.add("recieve")
      msgElement.textContent = currentChat.username + ":\n" + message.content + "\n"
    }
    msgElement.appendChild(dateElement)
    return msgElement
  
  }

  const throttle = <F extends (...args: Parameters<F>) => ReturnType<F>>( 
    fn: F, 
    delay: number, 
    options: { leading?: boolean; trailing?: boolean } = { 
      leading: true, 
      trailing: true, 
    } 
  ) => { 
    let timer: ReturnType<typeof setTimeout> | undefined 
   
    return (...args: Parameters<F>): void => { 
      if (timer !== undefined) return 
   
      options.leading && fn(...args) 
   
      timer = setTimeout(() => { 
        timer = undefined 
        options.trailing && fn(...args) 
      }, delay) 
    } 
  } 
  
  // lazyLoading of messages
  export const lazyLoading = throttle((chatId: number): void => {
    if (currentChat.chatId !== chatId) return
    const chatMessages = document.getElementById(
      `Chat${chatId}`
    ) as HTMLDivElement
  
    const buffer = messages.list.splice(currentChat.range.min, currentChat.range.max)
    for (const message of Object.values(buffer)) {
      chatMessages.append(createMessage(message))
    }
    if (currentChat.range.min + 10 > messages.list.length) {
      return
    }
    // Updates the range of messages that should be displayed
    currentChat.range.min += 10
    currentChat.range.max += 10
  
  }, 100)

  // Hides the chat and resets the state of the current chat
export const hideChat = (): void => {
    document.getElementById("chat-area")?.replaceChildren()
    Object.assign(currentChat, {
      username: "",
      userId: -1,
      chatId: -1,
      range: {
        min: 0,
        max: 10
      }
    })
  }

// toggleOnlineSection just shows/hides the online section
// Visible by default
export const toggleOnlineSection = (): void => {
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
  
  // toggleOfflineSection just shows/hides the offline section
  // Hidden by default
  export const toggleOfflineSection = (): void => {
    const offlineToggle = document.getElementById("offline-toggle") as HTMLElement
    const offlineUsers = document.getElementById(
      "offline-users"
    ) as HTMLUListElement
  
    if (offlineToggle.className === "bx bx-chevron-down") {
      offlineToggle.className = "bx bx-chevron-right"
      offlineUsers.style.display = "none"
    } else {
      offlineToggle.className = "bx bx-chevron-down"
      offlineUsers.style.display = "block"
    }
  }