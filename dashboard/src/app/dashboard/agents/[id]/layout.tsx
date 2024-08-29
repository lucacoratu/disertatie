

export default function Layout({
    children,
    metrics,
    logs,
    modal,
  }: {
    children: React.ReactNode,
    metrics: React.ReactNode,
    logs: React.ReactNode,
    modal: React.ReactNode,
  }) {
    return (
        <main className="h-full w-full flex flex-col gap-2 px-4">
          {children}
          {modal}            
          {metrics}
          {logs}
        </main>
    )
}