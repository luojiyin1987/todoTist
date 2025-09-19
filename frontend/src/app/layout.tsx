import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "TodoTist - Todo List App",
  description: "A simple todo list application built with Go + ConnectRPC and Next.js",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
