 


import { displayPosts } from "./posts.js"

// Category state
export const category = {
  selected: "all",
}

// Helper function to select functions. It affects the category state
export const categoriesSelector = (): void => {
  const categories = document.querySelectorAll(
    ".category-title"
  ) as NodeListOf<Element>

  // Selects the category
  categories.forEach((element) => {
    element.addEventListener("click", () => {
      if (element.classList.contains("selected")) {
        element.classList.remove("selected")
        category.selected = "all"
      } else {
        // De-select all other categories
        categories.forEach((otherElement) => {
          if (otherElement !== element) {
            otherElement.classList.remove("selected")
          }
        })
        element.classList.add("selected")
        category.selected = element.id
      }
      categoriesController()
    })
  })
  return
}

// Sends a request to the server based on selected category
const categoriesController = (): void => {
  switch (category.selected) {
    case "category-facts":
      displayPosts(`/api/posts/categories/facts`)
      break
    case "category-other":
      displayPosts("/api/posts/categories/other")
      break
    case "category-rumors":
      displayPosts("/api/posts/categories/rumors")
      break
    default:
      displayPosts("/api/posts")
      break
  }
  return
}
