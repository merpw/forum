// @type {import('next').NextConfig}
const nextConfig = {
  reactStrictMode: true,
  // swcMinify: true,
  images: {
    unoptimized: true,
  },
  // experimental: {
  //   scrollRestoration: true,
  // },
  async rewrites() {
    if (process.env.DEV_FORUM_BACKEND_REWRITE_URL) {
      return [
        {
          source: "/api/:path*",
          destination: `${process.env.DEV_FORUM_BACKEND_REWRITE_URL}/api/:path*`,
        },
      ]
    }
    return []
  },
  output: "standalone",
}
if (!process.env.FORUM_BACKEND_PRIVATE_URL) {
  console.warn("WARNING: FORUM_BACKEND_PRIVATE_URL is not defined.")
}
module.exports = nextConfig
