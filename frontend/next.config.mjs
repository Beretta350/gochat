/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",

  async rewrites() {
    // Em produção, usar variável de ambiente ou mesmo domínio
    const backendUrl = process.env.BACKEND_URL || "http://localhost:8080";

    return [
      // Proxy para API
      {
        source: "/api/:path*",
        destination: `${backendUrl}/api/:path*`,
      },
      // Proxy para WebSocket
      {
        source: "/ws",
        destination: `${backendUrl}/ws`,
      },
    ];
  },
};

export default nextConfig;
