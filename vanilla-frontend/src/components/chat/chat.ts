import { client } from "../../main.js"
import { Message } from "../../types"
import { state } from "../authorization/auth.js"
import { createElement } from "../utils.js"
import { renderChatList, renderChatMessages } from "./helpers/events.js"
import { sendWsObject } from "./helpers/sendobject.js"
import throttle from "./helpers/throttle.js"

export class Chat {
  constructor(
    readonly name: string,
    readonly userId: number,
    readonly chatId: number,
    public lastMessageId: number,
    public online: boolean
  ) {}

  messageIds: number[] = []
  unreadMessages = 0

  private range = 0 // Determines how many messages should be loaded

  typing = false

  // DOM for chat window
  chatArea = document.getElementById("chat-area") as HTMLElement
  chatWindow = createElement("div", "chat show-chat") as HTMLDivElement

  chatHeader = createElement("div", "chat-header") as HTMLDivElement
  chatName = createElement(
    "div",
    "chat-name",
    null,
    this.name
  ) as HTMLDivElement
  closeBtn = createElement(
    "a",
    "closebtn",
    "chat-window-close",
    null,
    "<i class='bx bx-x'>"
  ) as HTMLAnchorElement

  chatMessages = createElement(
    "div",
    "chat-messages",
    `chat${this.chatId}`
  ) as HTMLDivElement

  chatFormContainer = createElement(
    "div",
    "chat-form-container"
  ) as HTMLDivElement
  chatForm = createElement("form", "chat-form") as HTMLFormElement
  messageField = createElement("input", null, "chat-text") as HTMLInputElement
  messageSend = createElement(
    "button",
    null,
    "chat-send",
    "Send"
  ) as HTMLButtonElement



  open() {
    client.activeChat = this

    this.range = 0
    this.getMessages()
    this.chatMessages.replaceChildren()

    this.closeBtn.addEventListener("click", (e) => {
      e.preventDefault()
      this.close()
    })

    this.chatHeader.append(this.chatName, this.closeBtn)

    this.chatMessages.addEventListener("renderChatMessages", (e) => {
      setTimeout(() => {
        this.renderMessages()
        e.preventDefault()
      }, 100)
    })

    this.chatMessages.addEventListener("renderNewMessages", (e) => {
      setTimeout(() => {
        this.renderNewMessage()
        e.preventDefault()
      }, 100)
    })

    this.chatMessages.addEventListener("scroll", this.lazyLoading)

    this.chatForm.setAttribute("autocomplete", "off")
    this.messageField.setAttribute("maxlength", "150")
    this.messageSend.style.marginLeft = "4px"

    this.messageSend.addEventListener("click", (e) => {
      const content = this.messageField.value.trim().slice(0, 150)
      if (content.length === 0) {
        e.preventDefault()
        return
      }

      this.sendMessage(content)
      this.messageField.value = ""
      this.chatMessages.scrollTop = 0
      e.preventDefault()
    })

    this.chatForm.append(this.messageField, this.messageSend)
    this.chatFormContainer.append(this.chatForm)
    this.chatWindow.append(
      this.chatHeader,
      this.chatMessages,
      this.chatFormContainer
    )
    this.chatArea.replaceChildren(this.chatWindow)

    this.unreadMessages = 0
    window.dispatchEvent(renderChatList)

    this.chatMessages.dispatchEvent(renderChatMessages)
  }

  close() {
    this.chatMessages.replaceChildren()
    this.chatArea.replaceChildren()
    this.range = 0
    client.activeChat = null
  }

  private getMessages() {
    for (const id of this.messageIds.slice(this.range, this.range + 10)) {
      if (client.messages.has(id)) return
      sendWsObject({
        type: "get",
        item: {
          url: `/message/${id}`,
        },
      })
    }
  }

  private lazyLoading = throttle(() => {
    const scrollableHeight =
      this.chatMessages.scrollHeight - this.chatMessages.clientHeight - 50

    if (Math.abs(this.chatMessages.scrollTop) >= scrollableHeight) {
      this.range =
        this.range + 10 < this.messageIds.length
          ? this.range + 10
          : this.messageIds.length
      this.getMessages()
      this.chatMessages.dispatchEvent(renderChatMessages)
    }

    if (this.chatMessages.scrollTop >= 0 && this.unreadMessages > 0) {
      this.unreadMessages = 0
      window.dispatchEvent(renderChatList)
    }
  }, 200)

  private renderNewMessage(): void {
    const [id] = this.messageIds
    const messageElement = this.chatMessages.querySelector(`#msgId${id}`)
    if (messageElement) return

    const message = client.messages.get(id)
    if (!message){
      return this.renderNewMessage()
    }

    this.chatMessages.prepend(
      this.createMessage(message)
    )

    this.range += 1
  }

  private renderMessages(): void {
    for (const id of this.messageIds.slice(this.range, this.range + 10)) {
      const messageElement = this.chatMessages.querySelector(`#msgId${id}`)
      if (messageElement) {
        continue
      }

      const message = client.messages.get(id)
      if (!message) {
        return this.renderMessages()
      }

      this.chatMessages.appendChild(this.createMessage(message))
    }
  }

  private createMessage(message: Message): HTMLDivElement {
    const msgElement = createElement(
      "div",
      "message",
      `msgId${message.id}`
    ) as HTMLDivElement

    const date = new Date(message.timestamp)
    const formatDate = date
      .toLocaleString("en-GB", { timeZone: "EET" })
      .slice(0, 10)

    const dateElement = createElement("p", "date", null, formatDate)

    if (message.authorId === state.me.id) {
      msgElement.classList.add("send")
      msgElement.textContent = "You:\n" + message.content + "\n"
    } else if (message.authorId === -1) {
      msgElement.classList.add("status")
      msgElement.textContent = message.content + "\n"
    } else {
      msgElement.classList.add("recieve")
      msgElement.textContent = this.name + ":\n" + message.content + "\n"
    }
    msgElement.appendChild(dateElement)

    return msgElement
  }

  private sendMessage(content: string) {
    sendWsObject({
      type: "post",
      item: {
        url: `/chat/${this.chatId}/message`,
        data: {
          content: content,
        },
      },
    })
  }
}
