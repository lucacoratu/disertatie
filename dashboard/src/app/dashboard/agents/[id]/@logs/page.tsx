import { constants } from "@/app/constants";
import { DataTable } from "@/components/table/data-table";
import { columns, LogColumn } from './columns';
import { cookies } from "next/headers";

async function getLogs(agentId: string): Promise<LogShortResponse> {
	const cookie = cookies().get('session');
	//Create the URL where the logs will be fetched from
	const URL = `${constants.apiBaseURL}/agents/${agentId}/logs`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load logs");
	}
	//Parse the json data
	const logsResponse: LogShortResponse = await res.json();
	return logsResponse;
}

async function getLogsElastic(agentId: string): Promise<LogsShortElasticResponse> {
	const cookie = cookies().get('session');
	//Create the URL where the logs will be fetched from
	const URL = `${constants.apiBaseURL}/agents/${agentId}/logs-elastic`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 0}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load logs");
	}
	//Parse the json data
	const logsResponse: LogsShortElasticResponse = await res.json();
	return logsResponse;
}

async function getLogsPaginated(agentId: string, nextPage: string | null ): Promise<LogShortResponse> {
	const cookie = cookies().get('session');
	//Create the URL where the logs will be fetched from
  	var URL = "";
 	 if(nextPage === null) {
    URL = `${constants.apiBaseURL}/agents/${agentId}/logs`;
 	 } else {
	  URL = `${constants.apiBaseURL}/agents/${agentId}/logs?page=${nextPage}`;
  	}
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load logs");
	}
	//Parse the json data
	const logsResponse: LogShortResponse = await res.json();
	return logsResponse;
}

async function getFindingsStringFormat(): Promise<FindingClassificationString[]> {
	const cookie = cookies().get('session');
	//Create the URL where the findings classfication in string format will be fetched from
	const URL = `${constants.apiBaseURL}/findings/string`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load findings string");
	}
	//Parse the json data
	const findingsStringResponse: FindingClassificationStringResponse = await res.json();
	return findingsStringResponse.findingsString;
}

export default async function TableLogsPage({ params }: { params: { id: string } }) {
    const agentId: string = params.id;
	if(!agentId) {
		return;
	}
    //Get the logs
	//const log: LogShortResponse = await getLogs(agentId);
	const log: LogsShortElasticResponse = await getLogsElastic(agentId);

	//Get the findings classifications in string format
	const findingsClassficationString : FindingClassificationString[] = await getFindingsStringFormat();

	//Create the structure for the table
	const tableData: LogColumn[] = log.logs?.map((l, index) => {
		//Convert timestamp to local date time
		const logDate = new Date(l.timestamp * 1000);
		const logCol: LogColumn = {
			id: l.id,
			agentId: agentId,
			remoteip: l.remoteIp,
			requestmethod: l.request_preview.split(" ")[0],
			url: l.request_preview.split(" ")[1],
			response: l.response_preview.split(" ").slice(1, Infinity).join(" "),
			timestamp: logDate.toLocaleString(),
			findings: l.findings,
			findingsClassificationString: findingsClassficationString,
			rulefindings: l.ruleFindings,
			nextPage: ""
			// nextPage: log.nextPage, 
		}
		//Return the new column
		return logCol;
	});
    
    return (
        <div>
            <DataTable columns={columns} data={tableData} title="Logs"/>
        </div>
    );
}