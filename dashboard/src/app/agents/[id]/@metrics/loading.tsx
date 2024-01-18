import LoadingCustomPieChart from "@/components/ui/loading-piechart";

export default async function LoadingMetrics() {
    return (
        <div className="flex flex-row gap-4 flex-wrap">
            <div className="w-1/5 min-w-fit grow">
                <LoadingCustomPieChart title="Methods Distribution" numberLabels={6}/>
            </div>
            <div className="w-1/5 min-w-fit grow">
                <LoadingCustomPieChart title="Methods Distribution" numberLabels={6}/>
            </div>
            <div className="w-1/5 min-w-fit grow">
                <LoadingCustomPieChart title="Status Codes Distribution" numberLabels={4}/>
            </div>
            <div className="w-1/5 min-w-fit grow">
                <LoadingCustomPieChart title="Methods Distribution" numberLabels={6}/>
            </div>
        </div>
    );
}