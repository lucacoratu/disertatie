import { FC } from "react";
import {Skeleton} from "@/components/ui/skeleton"

type LoadingCustomPieChartProps = {
    numberLabels: number,
    title: string,
}

const LoadingCustomPieChart: FC<LoadingCustomPieChartProps> = ({ numberLabels, title }): JSX.Element => {
    return (
        <div className="py-6 flex flex-col gap-3 items-center dark:border-darksurface-400 border-2 dark:bg-darksurface-100 b-2 rounded-xl">
            <h3 className="font-bold text-sm">
                {title}
            </h3>
            <div className="flex flex-row items-center gap-6">
                <Skeleton className="h-60 w-60 rounded-full" />
                <div className="flex flex-col gap-2">
                    {
                        [...Array(numberLabels)].map((value) => (
                            <div key={value} className="flex flex-row gap-2">
                                <Skeleton className="h-[20px] w-[20px] rounded-none" />
                                <Skeleton className="h-[20px] w-[70px]" />
                            </div>
                        ))
                    }
                </div>
            </div>
        </div>
    );
}

export default LoadingCustomPieChart;