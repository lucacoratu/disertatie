export default function Layout({
    children,
  }: {
    children: React.ReactNode,
  }) {
    return (
        <main className="h-full w-full flex flex-col gap-4">
            {children}
        </main>
    )
}