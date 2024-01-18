import { constants } from "@/app/constants";
import CustomPieChart from "@/components/ui/piechart";
import CustomBarChart from "@/components/ui/barchart";

async function getMethodsMetrics(agentId: string) : Promise<MethodsMetrics[]>{
	//Create the URL where the metrics will be fetched from
	const URL = `${constants.apiBaseURL}/agents/${agentId}/logs-methods-metrics`;
	//Fetch the data (revalidate data after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load methods metrics");
	}
	//Parse the json data
	const metricsResponse: MethodMetricsResponse = await res.json();
	//Return data
	return metricsResponse.metrics;
}

async function getDaysMetrics(agentId: string) : Promise<DaysMetrics[]> {
	//Create the URL where the metrics will be fetched from
	const URL = `${constants.apiBaseURL}/agents/${agentId}/logs-each-day-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: { revalidate: 600}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load days metrics");
	}
	//Parse the json data
	const daysMetricsResponse: DaysMetricsResponse = await res.json();
	return daysMetricsResponse.metrics;
}

async function getStatusCodeMetrics(agentId: string) : Promise<StatusCodeMetrics[]> {
	//Create the URL where the metrics will be fetched from
	const URL = `${constants.apiBaseURL}/agents/${agentId}/logs-statuscode-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: { revalidate: 600}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load status code metrics");
	}
	//Parse the json data
	const statusCodeMetricsResponse: StatusCodeMetricsResponse = await res.json();
	return statusCodeMetricsResponse.metrics;
}

async function getIPAddressesMetrics(agentId: string) : Promise<IPMetrics[]> {
	//Create the URL where the metrics will be fetched from
	const URL = `${constants.apiBaseURL}/agents/${agentId}/logs-ipaddresses-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: { revalidate: 600}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load status code metrics");
	}
	//Parse the json data
	const ipAddressMetrics: IPAddressesMetricsResponse = await res.json();
	return ipAddressMetrics.metrics;
}

export default async function MetricsPage({ params }: { params: { id: string } }) {
    const agentId: string = params.id; 

	//Get metrics from the server for the pie chart
	const methodsMetrics: MethodsMetrics[] = await getMethodsMetrics(agentId);

	//Get the metrics from the server for the bar chart
	const daysMetrics: DaysMetrics[] = await getDaysMetrics(agentId);

	//Get the status code metrics
	const statusCodeMetrics: StatusCodeMetrics[] = await getStatusCodeMetrics(agentId);

	//Get the IP addresses metrics
	const ipAddressMetrics: IPMetrics[] = await getIPAddressesMetrics(agentId);
    
    return (
        <div className="flex flex-row gap-4 flex-wrap">
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
        </div>
    );
}