import { Configuration, OpenAIApi } from "openai"
import { AxiosError } from "axios"
import { cookies } from "next/headers"
import { serialize } from "cookie"

const configuration = new Configuration({
  apiKey: process.env.OPENAI_API_KEY,
})
const openai = new OpenAIApi(configuration)

// axios, which is used by OpenAIApi is not available in the Edge runtime
export const runtime = "nodejs"

// TODO: consider using websockets and moving this to separate service

export const POST = async (req: Request) => {
  if (!process.env.OPENAI_API_KEY) {
    return new Response("Temporarily unavailable", { status: 503 })
  }

  const tokenCookie = cookies().get("forum-token")

  if (!tokenCookie) {
    return new Response("Unauthorized", { status: 401 })
  }

  const response = await fetch(
    `${process.env.FORUM_BACKEND_PRIVATE_URL}/api/internal/check-session?token=${tokenCookie.value}`,
    {
      headers: {
        "Internal-Auth": process.env.FORUM_BACKEND_SECRET || "secret",
      },
      cache: "no-cache",
    }
  )

  if (!response.ok) {
    return new Response("Temporarily unavailable", { status: 503 })
  }

  const body = await response.json()
  if (body.error) {
    const cookie = serialize("forum-token", "", {
      expires: new Date(0),
    })
    return new Response("Unauthorized", { status: 401, headers: { "Set-Cookie": cookie } })
  }

  const requestData = await req.json()

  switch (requestData.action) {
    case "GENERATE_DESCRIPTION":
      return generateDescription(requestData.body)
    default:
      return new Response("Bad request", { status: 400 })
  }
}

const generateDescription = async (body: { title?: string; content?: string }) => {
  if (
    !body ||
    typeof body !== "object" ||
    typeof body.title !== "string" ||
    typeof body.content !== "string" ||
    body.title.length === 0 ||
    body.content.length === 0
  ) {
    return new Response("Invalid request", { status: 400 })
  }

  if (body.title.length > 50) {
    return new Response("Title is too long", { status: 400 })
  }

  if (body.content.length > 1000) {
    body.content = body.content.slice(0, 1000)
  }

  try {
    const result = await openai.createChatCompletion({
      model: "gpt-3.5-turbo",
      messages: [
        {
          role: "user",
          content: `
    Post title: <<${body.title}>>\n
    Post content: <<${body.content}>>\n
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
      return new Response("Temporarily unavailable", { status: 503 })
    }

    return new Response(responseMessage, { status: 200 })
  } catch (err) {
    if ((err as AxiosError).response?.status === 429) {
      return new Response("Too many requests", { status: 429 })
    }
    console.error(err)
    return new Response("Temporarily unavailable", { status: 503 })
  }
}
