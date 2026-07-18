import { Link, createFileRoute } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/")({ component: PublicPage });

function PublicPage() {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-3xl flex-col justify-center gap-6 px-6 py-24">
      <div className="flex flex-col gap-2">
        <h1 className="text-3xl font-semibold tracking-tight">Application workspace</h1>
        <p className="text-muted-foreground">Public entry point for your next product.</p>
      </div>
      <div>
        <Button asChild>
          <Link to="/dashboard">Open dashboard</Link>
        </Button>
      </div>
    </main>
  );
}
