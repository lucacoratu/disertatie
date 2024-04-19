import { constants } from "@/app/constants";
import FindingCard from "@/components/FindingsCard";
import RuleFindingCard from "@/components/RuleFindingCard";
import HTTPHighlighter from "@/components/HTTPHighlighter";
import ExploitCard from "@/components/ExploitCard";

async function GetLogFull(logid: string) : Promise<LogFull> {
    //Create the URL where the log will be fetched from
	const URL = `${constants.apiBaseURL}/logs/${logid}`;
    //Revalidate the data once every 10 mins
    const res = await fetch(URL, { cache: "no-store" });
    //Check if there was an error
    if(!res.ok) {
        console.log(res.ok);
        throw new Error("could not get log");
    }
    const logResponse: LogFullResponse = await res.json();
    //Return the data
    return logResponse.log;
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

async function getLogExploit(log_id: string): Promise<Exploit> {
    //Create the URL where the findings classfication in string format will be fetched from
	const URL = `${constants.apiBaseURL}/logs/${log_id}/exploit`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {next: {revalidate: 600}});
    //Check if an error occured
	if(!res.ok) {
		throw new Error("could not load log exploit");
	}
    const exploit: Exploit = await res.json();
    return exploit;
}

export default async function LogDetails({ params }: { params: { logid: string } }) {
    const logId: string = params.logid;

    //Get the full log from the API
    const fullLog: LogFull = await GetLogFull(logId);

    //Get the list of string descriptions of the findings based on the classification
    const findingsClassficationString : FindingClassificationString[] = await getFindingsStringFormat();

    //Get the log exploit
    const logExploit: Exploit = await getLogExploit(logId);

    return (
        <>
            <HTTPHighlighter log={fullLog}/>
            <div className="grid grid-cols-2 gap-4">
                <div className="flex flex-wrap flex-col gap-0 grow min-h-[200px] bg-card justify-center p-4 border b-2 rounded-xl">
                    <h2 className="text-xl">Findings</h2>
                    <FindingCard findings={fullLog.findings} findingsClassificationString={findingsClassficationString}/>
                </div>
                <div className="flex flex-wrap flex-col gap-0 grow min-h-[200px] bg-card justify-center p-4 border b-2 rounded-xl">
                    <h2 className="text-xl">Rule Findings</h2>
                    <RuleFindingCard findings={fullLog.ruleFindings}/>
                </div>
            </div>
            <ExploitCard exploit={logExploit}/>
        </>
    )
}