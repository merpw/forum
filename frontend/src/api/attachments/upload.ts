import axios from "axios"

const dummyResponse =
  "https://img.freepik.com/free-photo/red-white-cat-i-white-studio_155003-13189.jpg"

const uploadFile = async (file: File): Promise<string> => {
  const formData = new FormData()
  formData.append("file", file)

  return Promise.resolve(dummyResponse)

  return axios
    .post<string>("/api/attachments/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    })
    .then((res) => res.data)
}

export default uploadFile
