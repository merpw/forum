const edgeFetcher = async <T = unknown>(url: string, options: RequestInit = {}) =>
  fetch(process.env.FORUM_BACKEND_PRIVATE_URL + url, {
    ...options,
    headers: {
      ...options.headers,
      "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
    },
  }).then((res) => {
    if (!res.ok) throw new Error("Network response was not ok")
    return res.json() as Promise<T>
  })

export default edgeFetcher
