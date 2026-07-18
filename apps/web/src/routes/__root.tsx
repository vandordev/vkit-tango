import { HeadContent, Link, Scripts, createRootRoute } from "@tanstack/react-router";
import type { ErrorComponentProps } from "@tanstack/react-router";
import type { ReactNode } from "react";

import { QueryProvider } from "@/components/query-provider";
import { Button } from "@/components/ui/button";
import appCss from "@/styles.css?url";

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "Application Workspace" },
      { name: "description", content: "A reusable TanStack Start application workspace" },
    ],
    links: [{ rel: "stylesheet", href: appCss }],
  }),
  notFoundComponent: NotFoundPage,
  errorComponent: RouteErrorPage,
  shellComponent: RootDocument,
});

function FallbackLayout({ children }: { children: ReactNode }) {
  return <main className="mx-auto flex min-h-screen w-full max-w-xl flex-col justify-center gap-6 px-6 py-24">{children}</main>;
}

function NotFoundPage() {
  return (
    <FallbackLayout>
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold tracking-tight">Page not found</h1>
        <p className="text-muted-foreground">The page you are looking for does not exist or has moved.</p>
      </div>
      <div className="flex flex-wrap gap-3">
        <Button asChild>
          <Link to="/">Back to home</Link>
        </Button>
      </div>
    </FallbackLayout>
  );
}

function RouteErrorPage({ error, reset }: ErrorComponentProps) {
  return (
    <FallbackLayout>
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold tracking-tight">Something went wrong</h1>
        <p className="text-muted-foreground">Please try again. If the problem continues, return to the home page.</p>
        {/* eslint-disable-next-line turbo/no-undeclared-env-vars -- Vite injects DEV at build time. */}
        {import.meta.env.DEV && error instanceof Error ? <p className="text-sm text-muted-foreground">{error.message}</p> : null}
      </div>
      <div className="flex flex-wrap gap-3">
        <Button onClick={reset}>Try again</Button>
        <Button asChild variant="outline">
          <Link to="/">Back to home</Link>
        </Button>
      </div>
    </FallbackLayout>
  );
}

function RootDocument({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body className="min-h-screen bg-background font-sans text-foreground antialiased">
        <QueryProvider>{children}</QueryProvider>
        <Scripts />
      </body>
    </html>
  );
}
