import { NextApiRequest, NextApiResponse } from "next"
import { serialize } from "cookie"
import { Configuration, OpenAIApi } from "openai"
import { AxiosError } from "axios"

const configuration = new Configuration({
  apiKey: process.env.OPENAI_API_KEY,
})
const openai = new OpenAIApi(configuration)

// TODO: consider using websockets and moving this to separate service

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (!process.env.OPENAI_API_KEY) {
    return res.status(503).end("Temporarily unavailable")
  }
  const token = req.cookies["forum-token"]
  if (!token) {
    return res.status(401).end("Unauthorized")
  }

  const response = await fetch(
    `${process.env.FORUM_BACKEND_PRIVATE_URL}/api/internal/check-session?token=${token}`
  )
  const body = await response.json()
  if (body.error) {
    const cookie = serialize("forum-token", "", {
      expires: new Date(0),
    })
    return res.setHeader("Set-Cookie", cookie).status(401).end("Unauthorized")
  }

  switch (req.body.action) {
    case "GENERATE_DESCRIPTION":
      return generateDescription(req, res)
    default:
      return res.status(400).end("Invalid request")
  }
}

const generateDescription = async (req: NextApiRequest, res: NextApiResponse) => {
  if (
    typeof req.body.body.title !== "string" ||
    typeof req.body.body.content !== "string" ||
    req.body.body.title.length === 0 ||
    req.body.body.content.length === 0
  ) {
    return res.status(400).end("Invalid request")
  }

  if (req.body.body.title.length > 50) {
    return res.status(400).end("Title is too long")
  }

  if (req.body.body.content.length > 1000) {
    req.body.body.content = req.body.body.content.slice(0, 1000)
  }

  try {
    const result = await openai.createChatCompletion({
      model: "gpt-3.5-turbo",
      messages: [
        {
          role: "user",
          content: `
    Post title: <<${req.body.body.title}>>\n
    Post content: <<${req.body.body.content}>>\n
    Write a short creative TL;DR description of the post without quotes: \n
    <<`,
        },
      ],
      max_tokens: 200,
      n: 1,
      stream: false,
      stop: ">>",
    })

    const responseMessage = result.data.choices[0].message?.content

    if (!responseMessage?.length) {
      console.log("OpenAI error", result.data)
      return res.status(503).send("Temporarily unavailable")
    }

    res.status(200).json(responseMessage)
  } catch (err) {
    if ((err as AxiosError).response?.status === 429) {
      return res.status(429).send("Too many requests")
    }
    console.error(err)
    return res.status(503).send("Temporarily unavailable")
  }
}
