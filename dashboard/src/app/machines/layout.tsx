export default function Layout({
    children,
  }: {
    children: React.ReactNode,
  }) {
    return (
        <main className="h-full w-full flex flex-col py-2 gap-2 px-2 dark:bg-darksurface-200">            
            {children}
        </main>
    )
}