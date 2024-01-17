export default function Layout({
    children,
    metrics,
    tablelogs,
  }: {
    children: React.ReactNode,
    metrics: React.ReactNode,
    tablelogs: React.ReactNode,
  }) {
    return (
        <main className="h-full w-full flex flex-col py-2 gap-2 px-2 dark:bg-darksurface-200">            
            {metrics}
            {tablelogs}
            {children}
        </main>
    )
}