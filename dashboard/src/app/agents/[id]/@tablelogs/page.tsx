import { constants } from "@/app/constants";
import { DataTable } from "./data-table";
import { columns, LogColumn } from './columns';

async function getLogs(agentId: string): Promise<LogShort[]> {
	//Create the URL where the logs will be fetched from
	const URL = `${constants.apiBaseURL}/agents/${agentId}/logs`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load logs");
	}
	//Parse the json data
	const logsResponse: LogShortResponse = await res.json();
	return logsResponse.logs;
}

async function getFindingsStringFormat(): Promise<FindingClassificationString[]> {
	//Create the URL where the findings classfication in string format will be fetched from
	const URL = `${constants.apiBaseURL}/findings/string`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}});
	//Check if an error occured
	if(!res.ok) {
		throw new Error("could not load logs");
	}
	//Parse the json data
	const findingsStringResponse: FindingClassificationStringResponse = await res.json();
	return findingsStringResponse.findingsString;
}

export default async function TableLogsPage({ params }: { params: { id: string } }) {
    const agentId: string = params.id; 
    //Get the logs
	const logs: LogShort[] = await getLogs(agentId);
	//Get the findings classifications in string format
	const findingsClassficationString : FindingClassificationString[] = await getFindingsStringFormat();

	//Create the structure for the table
	const tableData: LogColumn[] = logs.map((log, index) => {
		//Convert timestamp to local date time
		const logDate = new Date(log.timestamp * 1000);
		const logCol: LogColumn = {
			id: log.id,
			agentId: agentId,
			remoteIp: log.remoteIp,
			requestMethod: log.request_preview.split(" ")[0],
			url: log.request_preview.split(" ")[1],
			response: log.response_preview.split(" ").slice(1, Infinity).join(" "),
			timestamp: logDate.toLocaleString(),
			findings: log.findings,
			findingsClassificationString: findingsClassficationString,
		}
		//Return the new column
		return logCol;
	})
    
    return (
        <div>
            <DataTable columns={columns} data={tableData}/>
        </div>
    );
}