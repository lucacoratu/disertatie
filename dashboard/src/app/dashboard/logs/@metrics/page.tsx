import { constants } from "@/app/constants";
import CustomPieChart from "@/components/ui/piechart";

async function GetAgentMetrics(): Promise<AgentMetrics[]> {
    //Create the URL where the logs will be fetched from
	const URL = `${constants.apiBaseURL}/logs/agent-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load agent metrics");
	}
	//Parse the json data
	const agentMetricsRes: AgentMetricsResponse = await res.json();
	return agentMetricsRes.metrics;
}

async function GetClassificationMetrics(): Promise<ClassificationMetrics> {
    //Create the URL where the logs will be fetched from
	const URL = `${constants.apiBaseURL}/logs/classification-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load classification metrics");
	}
	//Parse the json data
	const agentMetricsRes: ClassificationMetricsResponse = await res.json();
	return agentMetricsRes.metrics;
}

export default async function MetricsPage() {
	// //Get metrics from the server for the pie chart
	// const methodsMetrics: MethodsMetrics[] = await getMethodsMetrics(agentId);

	// //Get the metrics from the server for the bar chart
	// const daysMetrics: DaysMetrics[] = await getDaysMetrics(agentId);

	// //Get the status code metrics
	// const statusCodeMetrics: StatusCodeMetrics[] = await getStatusCodeMetrics(agentId);

	// //Get the IP addresses metrics
	// const ipAddressMetrics: IPMetrics[] = await getIPAddressesMetrics(agentId);
    
    //Get the agent metrics from the API
    const agentMetrics: AgentMetrics[] = await GetAgentMetrics();

    //Get the classification metrics from the API
    const classificationMetrics: ClassificationMetrics = await GetClassificationMetrics();

    console.log(classificationMetrics);

    return (
        <div className="flex flex-row gap-4 flex-wrap">
            <div className="w-1/5 min-w-fit grow">
                <CustomPieChart labels={agentMetrics.map(({agentName}) => agentName)} values={agentMetrics.map(({count}) => count)} title="Agents Log Counts"/>
            </div>
            {/* <div className="w-1/5 min-w-fit grow">
                <CustomBarChart labels={daysMetrics.map(({date}) => date)} values={daysMetrics.map(({count}) => count)} title="No. Requests Over Time"/>
            </div> */}
            <div className="w-1/5 min-w-fit grow">
                <CustomPieChart labels={['Classified', 'Unclassified']} values={[classificationMetrics.classifiedCount, classificationMetrics.unclassifiedCount]} title="Status Codes Distribution"/>
            </div>
            {/* <div className="w-1/5 min-w-fit grow">
                <CustomBarChart labels={ipAddressMetrics.map(({ipAddress}) => ipAddress)} values={ipAddressMetrics.map(({count}) => count)} title="Remote IP Addresses"/>
            </div> */}
        </div>
    );
}