export default function Layout({
    children,
    tablealllogs,
  }: {
    children: React.ReactNode,
    tablealllogs: React.ReactNode,
  }) {
    return (
        <main className="w-full flex flex-col gap-4 p-4 md:px-8">            
            {children}
            {tablealllogs}
        </main>
    )
}