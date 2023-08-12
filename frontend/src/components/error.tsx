import { AnimatePresence, motion } from "framer-motion"
import { FC } from "react"

export const FormError: FC<{ error: string | null }> = ({ error }) => (
  <AnimatePresence>
    {error && (
      <motion.div
        className={
          "transition ease-in-out -translate-y-1 px-3 py-2 my-1 text-sm bg-error text-error-content rounded-lg"
        }
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        role={"alert"}
      >
        <span className={"font-medium"}>{error}</span>
      </motion.div>
    )}
  </AnimatePresence>
)
