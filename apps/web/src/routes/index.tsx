import { Link, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({ component: PublicPage });

function PublicPage() {
  return (
    <main>
      <h1>Application workspace</h1>
      <p>Public entry point for your next product.</p>
      <Link to="/dashboard">Open dashboard</Link>
    </main>
  );
}
