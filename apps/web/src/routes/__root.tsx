import { Box, Button, ColorSchemeScript, createTheme, Group, MantineProvider, Stack, Text, Title } from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { HeadContent, Scripts, createRootRoute, useNavigate } from "@tanstack/react-router";
import type { ErrorComponentProps } from "@tanstack/react-router";
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
  notFoundComponent: NotFoundPage,
  errorComponent: RouteErrorPage,
  shellComponent: RootDocument,
});

function FallbackLayout({ children }: { children: ReactNode }) {
  return (
    <Box component="main" maw={560} mx="auto" px="md" py={96}>
      <Stack gap="lg">{children}</Stack>
    </Box>
  );
}

function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <FallbackLayout>
      <Stack gap="xs">
        <Title order={1}>Page not found</Title>
        <Text c="dimmed">The page you are looking for does not exist or has moved.</Text>
      </Stack>
      <Group>
        <Button onClick={() => void navigate({ to: "/" })}>Back to home</Button>
      </Group>
    </FallbackLayout>
  );
}

function RouteErrorPage({ error, reset }: ErrorComponentProps) {
  const navigate = useNavigate();

  return (
    <FallbackLayout>
      <Stack gap="xs">
        <Title order={1}>Something went wrong</Title>
        <Text c="dimmed">Please try again. If the problem continues, return to the home page.</Text>
        {/* eslint-disable-next-line turbo/no-undeclared-env-vars -- Vite injects DEV at build time. */}
        {import.meta.env.DEV && error instanceof Error ? <Text c="dimmed" size="sm">{error.message}</Text> : null}
      </Stack>
      <Group>
        <Button onClick={reset}>Try again</Button>
        <Button onClick={() => void navigate({ to: "/" })} variant="default">
          Back to home
        </Button>
      </Group>
    </FallbackLayout>
  );
}

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
