export const wrapWith = (textAreaRef: HTMLTextAreaElement, before: string, after: string) => {
  let { selectionStart, selectionEnd } = textAreaRef

  if (selectionStart === selectionEnd) {
    // select word
    while (selectionStart > 0 && textAreaRef.value[selectionStart - 1] !== " ") {
      selectionStart--
    }
    while (selectionEnd < textAreaRef.value.length && textAreaRef.value[selectionEnd] !== " ") {
      selectionEnd++
    }
  }

  const beforeText = textAreaRef.value.slice(0, selectionStart)
  const afterText = textAreaRef.value.slice(selectionEnd)
  const selected = textAreaRef.value.slice(selectionStart, selectionEnd)

  textAreaRef.value = beforeText + before + selected + after + afterText

  textAreaRef.selectionStart = textAreaRef.selectionEnd =
    beforeText.length + selected.length + before.length

  textAreaRef.focus()
}
