import Sidebar from "@/components/Sidebar";

export const metadata = {
  title: "Dashboard — Aura Amber",
};

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <main className="flex-1 px-6 py-8 md:px-10">{children}</main>
    </div>
  );
}
