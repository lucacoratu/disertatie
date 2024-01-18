import { constants } from "@/app/constants";
import HTTPHighlighter from "@/components/HTTPHighlighter";

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

export default async function LogDetails({ params }: { params: { logid: string } }) {
    const logId: string = params.logid;

    //Get the full log from the API
    const fullLog: LogFull = await GetLogFull(logId);

    return (
        <HTTPHighlighter log={fullLog}/>
    )
}