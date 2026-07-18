import { ColorSchemeScript, createTheme, MantineProvider } from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { HeadContent, Scripts, createRootRoute } from "@tanstack/react-router";
import type { ReactNode } from "react";

import { QueryProvider } from "@/components/query-provider";
import appCss from "@/styles.css?url";

const theme = createTheme({
  primaryColor: "oriskin",
  defaultRadius: "md",
  fontFamily: '"Space Grotesk", system-ui, sans-serif',
  headings: { fontFamily: '"Space Grotesk", system-ui, sans-serif' },
  colors: {
    oriskin: ["#fff1f0", "#ffe1df", "#ffc4bf", "#ff9d96", "#f87168", "#e84f45", "#d93a30", "#b82d25", "#982821", "#7e251f"],
  },
});

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "Application Workspace" },
      { name: "description", content: "A reusable TanStack Start application workspace" },
    ],
    links: [
      { rel: "preconnect", href: "https://fonts.googleapis.com" },
      { rel: "preconnect", href: "https://fonts.gstatic.com", crossOrigin: "anonymous" },
      { rel: "stylesheet", href: "https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300..700&display=swap" },
      { rel: "stylesheet", href: appCss },
    ],
  }),
  shellComponent: RootDocument,
});

function RootDocument({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <head>
        <ColorSchemeScript defaultColorScheme="light" />
        <HeadContent />
      </head>
      <body>
        <MantineProvider defaultColorScheme="light" theme={theme}>
          <QueryProvider>
            <Notifications position="top-right" />
            {children}
          </QueryProvider>
        </MantineProvider>
        <Scripts />
      </body>
    </html>
  );
}
