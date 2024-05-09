export default function Layout({
    children,
    tableclassifiedlogs,
  }: {
    children: React.ReactNode,
    tableclassifiedlogs: React.ReactNode,
  }) {
    return (
        <main className="w-full flex flex-col gap-4 p-4 md:px-8">            
            {children}
            {tableclassifiedlogs}
        </main>
    )
}