import { backendUrl } from "../main.js"
import { ChatUser } from "../types"
import { sendObject } from "./ws.js"
import { userInfo } from "./auth.js"
// This file is dedicated to sorting and displaying chat users in the sidebar.
export const chatList = {
  ids: [] as number[],
  chatIds: new Map<number, number>, //userId, chatId
  users: [] as ChatUser[]
}

export const messages = {
  list: [] as unknown[]
}

async function getUserIds() {
  await fetch(`${backendUrl}/api/users`)
  .then((r) => r.json())
  .then((data) => {
    for (const id of data) {
      chatList.ids.push(id)
      } 
  })
}

async function getChatUsers() {
  Object.assign(chatList, {users: [], ids: []})
  await getUserIds()

  for (const id of chatList.ids){
    sendObject({
        type: "post",
        item: {
          url: "/chat/create",
          data: {
            userId: id
          }
        }
      })

    await fetch(`${backendUrl}/api/user/${id}`)
    .then((r) => r.json())
    .then((data) => {
      if (data.id !== userInfo.Id){
        const user: ChatUser = {
          Id: data.id,
          Name: data.name,
          UnreadMsg: false,
          Online: false
        }
        chatList.users.push(user)
      }
    }
    )
  }

  //   ws.send(JSON.stringify({
  //       type: "get",
  //       item: {
  //         url: `/users/online`,
  //         data: {
  //           content: message,
  //         }
  //       }
  //     })
  // ) TODO: Fix so that online users and UnreadMsg are updated in chatList
}

const showChat = async (id: number) => {
  const chatId = chatList.chatIds.get(id) as number
  const chatArea = document.createElement("div")

  const chat = document.createElement("div")
  chat.className = "chat show-chat"

  const chatMessages = document.createElement("div")
  chatMessages.className = "chat-messages"
  chatMessages.id = `Chat${chatId}`
  getMessages(chatId)

  const chatFormContainer = document.createElement("div")
  chatFormContainer.className = "chat-form-container"

  const chatForm = document.createElement("form")
  chatForm.className ="chat-form"
  
  chatArea.replaceChildren(chat)
}
//   <div id="chat-test">
//     <div id="chat-messages">
//       <div class="message send">TEST TEST TEST 123</div>
//       <div class="message recieve"></div>
//     </div>
//     <div class="chat-form-container">
//       <form id="chat-form">
//       <input type="text" id="chat-text">
//       <input type="submit" id="chat-send" value="Send">
//       </form>
//     </div>
// </div>

const getMessages = async (chatId: number) => {
  sendObject(
    {
      type: "get",
      item: {
        url: `/chat/${chatId}/messages`
      }
    })
}

export const displayChatUsers = async () => {
  const onlineList = document.getElementById("online-users") as HTMLUListElement,
  offlineList = document.getElementById("offline-users") as HTMLUListElement,
  onlineTitle = document.getElementById("online-title") as HTMLElement,
  offlineTitle = document.getElementById("offline-title") as HTMLElement

  /* Add eventlisteners to show/hide users in chat list */
  onlineTitle.addEventListener("click", toggleOnline)
  offlineTitle.addEventListener("click", toggleOffline)

  await getChatUsers()

  chatList.users.sort((a, b) => {
    const name1 = a.Name.toLowerCase()
    const name2 = b.Name.toLowerCase()
    if (name1 < name2) {
      return -1
    }
    if (name1 > name2) {
      return 1
    }
    return 0
  })

  for (const u of chatList.users) {
    const user = document.createElement("li")
    const name = document.createElement("p")
    name.textContent = `${u.Name} `

    user.appendChild(name)

    if (u.UnreadMsg) {
      const unread = document.createElement("i")
      unread.className = "bx bx-message-dots"
      user.appendChild(unread)
    }

    user.addEventListener("click", () => {
      showChat(u.Id)
    })

    if (u.Online) {
      user.className = "online"
      onlineList.appendChild(user)
    } else {
      user.className = "offline"
      offlineList.appendChild(user)
    }
  }
  sendObject(
    {
      type: "get",
      item: {
        url: "/chat/all"
      }
    }
  )
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
  const offlineUsers = document.getElementById("offline-users") as HTMLUListElement
  if (!offlineToggle) return

  if (offlineToggle.className === "bx bx-chevron-down") {
    offlineToggle.className = "bx bx-chevron-right"
    offlineUsers.style.display = "none"
  } else {
    offlineToggle.className = "bx bx-chevron-down"
    offlineUsers.style.display = "block"
  }
}
