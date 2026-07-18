import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/dashboard")({ component: DashboardPage });

function DashboardPage() {
  return (
    <main>
      <h1>Dashboard</h1>
      <p>Authenticated access can be added by the product built from this template.</p>
    </main>
  );
}
