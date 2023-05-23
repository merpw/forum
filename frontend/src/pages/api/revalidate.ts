import { NextApiRequest, NextApiResponse } from "next"

// TODO: revalidate app router pages

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
    if (!url) {
      return res.status(400).send("No url provided")
    }

    return res.json({ revalidated: true })
  } catch (err) {
    if ((err as Error).message.includes("Invalid response")) {
      return res.status(400).send((err as Error).message)
    }
    console.error(err)
    return res.status(500).send("Error revalidating")
  }
}
