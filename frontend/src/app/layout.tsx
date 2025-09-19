import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "TodoTist - Todo List App",
  description: "A simple todo list application built with Go + ConnectRPC and Next.js",
};

/**
 * Root layout component that wraps pages with the base HTML structure.
 *
 * Renders the <html lang="en"> element and a <body> with the "antialiased"
 * class, placing the provided `children` inside the body.
 *
 * @param children - React nodes to be rendered inside the document body.
 * @returns The root JSX layout for the application.
 */
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
