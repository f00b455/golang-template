/** @type {import('next').NextConfig} */
const nextConfig = {
  transpilePackages: ['@ts-template/shared'],
  output: 'export',
  trailingSlash: true,
  images: {
    unoptimized: true
  }
}

module.exports = nextConfig