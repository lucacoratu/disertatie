export default function Layout({
    children,
    tableclassifiedlogs,
    metrics,
  }: {
    children: React.ReactNode,
    tableclassifiedlogs: React.ReactNode,
    metrics: React.ReactNode,
  }) {
    return (
        <main className="w-full flex flex-col gap-4">            
            {children}
            {metrics}
            {tableclassifiedlogs}
        </main>
    )
}