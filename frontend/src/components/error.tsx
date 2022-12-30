import { AnimatePresence, motion } from "framer-motion"
import { FC } from "react"

export const FormError: FC<{ error: string | null }> = ({ error }) => (
  <AnimatePresence>
    {error && (
      <motion.div
        className={
          "transition ease-in-out -translate-y-1 p-4 mb-4 text-sm text-red-700 bg-red-100 rounded-lg dark:bg-inherit dark:border-2 dark:border-red-900 dark:text-white"
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
