/**
 * @type {import('next').NextConfig}
 */
const nextConfig = {
  output: "export",
  skipTrailingSlashRedirect: true,
  distDir: "dist",
  compress: false,
};

module.exports = nextConfig;
