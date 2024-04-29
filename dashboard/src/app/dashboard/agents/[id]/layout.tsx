

export default function Layout({
    children,
    metrics,
    tablelogs,
    modal,
  }: {
    children: React.ReactNode,
    metrics: React.ReactNode,
    tablelogs: React.ReactNode,
    modal: React.ReactNode,
  }) {
    return (
        <main className="h-full w-full flex flex-col py-4 gap-2 px-4">
          {children}
          {modal}            
          {metrics}
          {tablelogs}
        </main>
    )
}