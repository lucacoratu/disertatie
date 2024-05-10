export default function Layout({
    children,
    metrics,
    tablealllogs,
  }: {
    children: React.ReactNode,
    metrics: React.ReactNode,
    tablealllogs: React.ReactNode,
  }) {
    return (
        <main className="w-full flex flex-col gap-4">            
            {children}
            {metrics}
            {tablealllogs}
        </main>
    )
}