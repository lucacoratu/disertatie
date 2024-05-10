import { Skeleton } from "@/components/ui/skeleton";

export default async function LoadingTableLogs() {
    const numberEntriesTable = 10;

    return (
        <div className="flex w-full flex-col">
            <div className="mt-6 flex flex-row justify-between">
                <Skeleton className="w-[200px] h-[40px] dark:bg-darksurface-100"/>
                <Skeleton className="w-[100px] h-[40px] dark:bg-darksurface-100"/>
            </div>
            <div className="mt-4 flex flex-col justify-center w-full gap-0 border-2 dark:border-darksurface-400">
                {[...Array(numberEntriesTable)].map((_, index) => {
                    return index % 2 == 0 ? <Skeleton key={index} className="h-[40px] dark:bg-darksurface-100"/> : <Skeleton key={index} className="h-[40px] dark:bg-darksurface-300"/>
                })}
            </div>
        </div>
    );
}