import { ScrollArea } from "@/components/ui/scroll-area"

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
        <main className="h-full w-full flex flex-col py-2 gap-2 px-2 dark:bg-darksurface-200">
          {children}
          {modal}            
          {metrics}
          {tablelogs}
        </main>
    )
}