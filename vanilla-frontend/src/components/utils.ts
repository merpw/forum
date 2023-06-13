export const createElement = (
  element: string,
  className: string | null = null,
  id: string | null = null,
  textContent: string | null = null,
  innerHtml: string | null = null
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
  if (innerHtml) {
    newElement.innerHTML = innerHtml
  }
  return newElement
}

export interface iterator {
  [key: string]: any
}
