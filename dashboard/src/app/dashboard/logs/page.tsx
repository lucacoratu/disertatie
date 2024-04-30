import CustomPieChart from "@/components/ui/piechart";
import CustomBarChart from "@/components/ui/barchart";

export default async function Logs({ params }: { params: { id: string } }) {
    return (
        <>
            {/* <div className="flex flex-row gap-4 flex-wrap">
                <div className="w-1/5 min-w-fit grow">
                    <CustomPieChart labels={methodsMetrics.map(({method}) => method)} values={methodsMetrics.map(({count}) => count)} title="Methods Distribution"/>
                </div>
                <div className="w-1/5 min-w-fit grow">
                    <CustomBarChart labels={daysMetrics.map(({date}) => date)} values={daysMetrics.map(({count}) => count)} title="No. Requests Over Time"/>
                </div>
                <div className="w-1/5 min-w-fit grow">
                    <CustomPieChart labels={statusCodeMetrics.map(({statusCode}) => statusCode)} values={statusCodeMetrics.map(({count}) => count)} title="Status Codes Distribution"/>
                </div>
                <div className="w-1/5 min-w-fit grow">
                    <CustomBarChart labels={ipAddressMetrics.map(({ipAddress}) => ipAddress)} values={ipAddressMetrics.map(({count}) => count)} title="Remote IP Addresses"/>
                </div>
            </div> */}
        </>
    )
}