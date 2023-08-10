import axios from "axios"

const uploadFile = async (file: File): Promise<string> => {
  const formData = new FormData()
  formData.append("file", file)

  return axios
    .post<string>("/api/attachments/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    })
    .then((res) => res.data)
}

export default uploadFile
