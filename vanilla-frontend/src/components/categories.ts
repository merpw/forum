import { displayPosts } from "./posts.js"

// Category state
export const category = {
  selected: "all",
}

// Helper function to select functions. It affects the category state
export const categoriesSelector = async () => {
  const categories = document.querySelectorAll(
    ".category-title"
  ) as NodeListOf<Element>

  categories.forEach((categoryElement) => {
    categoryElement.addEventListener("click", () => {
      if (categoryElement.classList.contains("selected")) {
        categoryElement.classList.remove("selected")
        category.selected = "all"
      } else {
        categories.forEach((otherCategoryElement) => {
          if (otherCategoryElement !== categoryElement) {
            otherCategoryElement.classList.remove("selected")
          }
        })
        categoryElement.classList.add("selected")
        category.selected = categoryElement.id
      }
      categoriesController()
    })
  })
  return
}

const categoriesController = () => {
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
