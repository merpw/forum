// Creates and returns a HTMLElement.
// Current options:
// className, id, textContent, innerHTML.
export const createElement = (
  element: string,
  className: string | null = null,
  id: string | null = null,
  textContent: string | null = null,
  innerHTML: string | null = null
): HTMLElement => {
  const newElement = document.createElement(element)
  if (className) {
    newElement.className = className
  }
  if (id) {
    newElement.id = id
  }
  if (textContent) {
    newElement.textContent = textContent
  }
  if (innerHTML) {
    newElement.innerHTML = innerHTML
  }
  return newElement
}

export interface iterator {
  [key: string]: any
}
