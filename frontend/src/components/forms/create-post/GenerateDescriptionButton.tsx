import { Dispatch, FC, SetStateAction, useState } from "react"

import { generateDescription } from "@/api/posts/create"

export const GenerateDescriptionButton: FC<{
  formRef: React.RefObject<HTMLFormElement>
  setFormError: Dispatch<SetStateAction<string | null>>
}> = ({ formRef, setFormError }) => {
  const [isLoading, setIsLoading] = useState(false)

  return (
    <button
      onClick={async () => {
        const formData = new FormData(formRef.current as HTMLFormElement)

        const formFields = {
          title: formData.get("title") as string,
          content: formData.get("content") as string,
        }

        if (!formFields.title) {
          return setFormError("Title is required")
        }
        if (!formFields.content) {
          return setFormError("Content is required")
        }
        setIsLoading(true)
        generateDescription(formFields)
          .then((description) => {
            setFormError(null)
            const descriptionTextarea = formRef.current?.querySelector(
              "[name=description]"
            ) as HTMLInputElement
            descriptionTextarea.value = description
          })
          .catch((err) => setFormError(err.message))
          .finally(() => setIsLoading(false))
      }}
      type={"button"}
      className={
        "btn btn-sm transition-none hover:opacity-100 hover:gradient-text hover:border-primary btn-outline font-normal mb-3 self-center font-xs"
      }
    >
      {isLoading ? (
        <span className={"text-primary loading loading-ring"} />
      ) : (
        <svg
          xmlns={"http://www.w3.org/2000/svg"}
          fill={"none"}
          viewBox={"0 0 24 24"}
          strokeWidth={1}
          stroke={"currentColor"}
          className={"w-5 h-5 mr-1 fill-primary"}
        >
          <path
            strokeLinecap={"round"}
            strokeLinejoin={"round"}
            d={
              "M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z"
            }
          />
        </svg>
      )}
      Generate description
    </button>
  )
}
