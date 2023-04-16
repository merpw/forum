import { NextApiRequest, NextApiResponse } from "next"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  // Secret token is optional, by default frontend /api/ is not public
  if (req.query.secret !== process.env.FRONTEND_REVALIDATE_TOKEN) {
    return res.status(401).json({ message: "Invalid token" })
  }
  if (req.method !== "POST") {
    return res.status(405).end()
  }

  try {
    const url = req.query.url as string
    await res.revalidate(url)
    return res.json({ revalidated: true })
  } catch (err) {
    if ((err as Error).message.includes("Invalid response 404")) {
      return res.status(400).send("Revalidated page is 404")
    }
    console.error(err)
    return res.status(500).send("Error revalidating")
  }
}
