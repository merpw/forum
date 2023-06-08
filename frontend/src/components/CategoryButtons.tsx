"use client"

import { FC } from "react"
import Link from "next/link"
import { useParams } from "next/navigation"

const CategoryButtons: FC<{ categories: string[] }> = ({ categories }) => {
  const { name } = useParams()

  const activeCategory = name ?? "NONE"

  return (
    <div className={"mb-5"}>
      <div className={"text-center space-x-1"}>
        <div className={"font-light mb-1 text-center text-info"}>Categories:</div>
        <span className={"font-black gradient-text"}>{"•"}</span>
        {categories.map((category, key) => (
          <span key={key} className={""}>
            <Link
              href={`/category/${category}`}
              className={
                "btn btn-xs font-light" +
                " " +
                (activeCategory.toLowerCase() === category.toLowerCase()
                  ? "btn-secondary font-normal"
                  : "btn-neutral")
              }
            >
              {category}
            </Link>
          </span>
        ))}
        <span className={"font-black gradient-text"}>{"•"}</span>
      </div>
    </div>
  )
}

export default CategoryButtons
