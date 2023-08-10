export const wrapWith = (
  textAreaRef: HTMLTextAreaElement,
  before: string,
  after: string,
  autoSelect: "word" | "line" = "word"
) => {
  let { selectionStart, selectionEnd } = textAreaRef

  const backupSelectionStart = selectionStart
  const backupSelectionEnd = selectionEnd

  if (selectionStart === selectionEnd) {
    if (autoSelect === "line") {
      // select line
      while (selectionStart > 0 && textAreaRef.value[selectionStart - 1] !== "\n") {
        selectionStart--
      }
      while (selectionEnd < textAreaRef.value.length && textAreaRef.value[selectionEnd] !== "\n") {
        selectionEnd++
      }
    }
    if (autoSelect === "word") {
      // select word
      while (selectionStart > 0 && textAreaRef.value[selectionStart - 1].match(/\S/)) {
        selectionStart--
      }
      while (
        selectionEnd < textAreaRef.value.length &&
        textAreaRef.value[selectionEnd].match(/\S/)
      ) {
        selectionEnd++
      }
    }

    textAreaRef.select()
    textAreaRef.setSelectionRange(selectionStart, selectionEnd)
  }

  const selectedText = textAreaRef.value.substring(selectionStart, selectionEnd)

  // This is deprecated but the only way to insert text at a specific position without breaking Undo
  // Still works in all browsers
  document.execCommand("insertText", false, before + selectedText + after)

  textAreaRef.setSelectionRange(
    backupSelectionStart + before.length,
    backupSelectionEnd + before.length
  )

  textAreaRef.focus()
}
