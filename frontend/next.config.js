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
    return [
      {
        source: "/api/:path*",
        destination: `http://${process.env.FORUM_BACKEND_LOCALHOST}/api/:path*`,
      },
    ]
  },
}

module.exports = nextConfig
