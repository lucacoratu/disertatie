import { constants } from "@/app/constants";
import CustomPieChart from "@/components/ui/piechart";
import CustomBarChart from "@/components/ui/barchart";
import { cookies } from "next/headers";

async function GetAgentMetrics(): Promise<AgentMetrics[]> {
	const cookie = cookies().get('session');
    //Create the URL where the logs will be fetched from
	const URL = `${constants.apiBaseURL}/logs/agent-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load agent metrics");
	}
	//Parse the json data
	const agentMetricsRes: AgentMetricsResponse = await res.json();
	return agentMetricsRes.metrics;
}

async function GetClassificationMetrics(): Promise<ClassificationMetrics> {
	const cookie = cookies().get('session');
    //Create the URL where the logs will be fetched from
	const URL = `${constants.apiBaseURL}/logs/classification-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load classification metrics");
	}
	//Parse the json data
	const agentMetricsRes: ClassificationMetricsResponse = await res.json();
	return agentMetricsRes.metrics;
}

async function getIPAddressesMetrics() : Promise<IPMetrics[]> {
	const cookie = cookies().get('session');
	//Create the URL where the metrics will be fetched from
	const URL = `${constants.apiBaseURL}/logs/ip-address-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: { revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load status code metrics");
	}
	//Parse the json data
	const ipAddressMetrics: IPAddressesMetricsResponse = await res.json();
	return ipAddressMetrics.metrics;
}

export default async function MetricsPage() {    
    //Get the agent metrics from the API
    const agentMetrics: AgentMetrics[] = await GetAgentMetrics();

    //Get the classification metrics from the API
    const classificationMetrics: ClassificationMetrics = await GetClassificationMetrics();

    //Get the ip addresses metrics from the API
	const ipAddressesMetrics: IPMetrics[] = await getIPAddressesMetrics();

    return (
        <div className="flex flex-row gap-4 flex-wrap">
            <div className="w-1/5 min-w-fit grow">
                <CustomPieChart labels={agentMetrics.map(({agentName}) => agentName)} values={agentMetrics.map(({count}) => count)} title="Agents Log Counts"/>
            </div>
			<div className="w-1/5 min-w-fit grow">
                <CustomBarChart labels={ipAddressesMetrics.map(({ipAddress}) => ipAddress)} values={ipAddressesMetrics.map(({count}) => count)} title="Remote IP Addresses"/>
            </div>
            {/* <div className="w-1/5 min-w-fit grow">
                <CustomBarChart labels={daysMetrics.map(({date}) => date)} values={daysMetrics.map(({count}) => count)} title="No. Requests Over Time"/>
            </div> */}
            <div className="w-1/5 min-w-fit grow">
                <CustomPieChart labels={['Classified', 'Unclassified']} values={[classificationMetrics.classifiedCount, classificationMetrics.unclassifiedCount]} title="Classification Counts"/>
            </div>

        </div>
    );
}