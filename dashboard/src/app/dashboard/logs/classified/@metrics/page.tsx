import { constants } from "@/app/constants";
import FindingsRadarChart from "@/components/FindingsRadarChart";
import CustomBarChart from "@/components/ui/barchart";
import { cookies } from "next/headers";

async function GetRuleFindingsMetrics(): Promise<FindingsMetrics[]> {
	const cookie = cookies().get('session');
    //Create the URL where the logs will be fetched from
    const URL = `${constants.apiBaseURL}/findings/rule/metrics`;
    //Fetch the data (revalidate after 10 minutes)
    const res = await fetch(URL, {cache: "no-store", headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
    //Check if an error occured
    if(!res.ok) {
      throw new Error("could not load rule findings metrics");
    }
    //Parse the json data
    const findingsMetrics: FindingsMetricsResponse = await res.json();
    return findingsMetrics.metrics;
}

async function GetRuleIdsMetrics(): Promise<FindingsMetrics[]> {
	const cookie = cookies().get('session');
	//Create the URL where the logs will be fetched from
	const URL = `${constants.apiBaseURL}/findings/rule/id-metrics`;
	//Fetch the data (revalidate after 10 minutes)
	const res = await fetch(URL, {cache: "no-store", headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
	//Check if an error occured
	if(!res.ok) {
	  throw new Error("could not load rule ids metrics");
	}
	//Parse the json data
	const findingsMetrics: FindingsMetricsResponse = await res.json();
	return findingsMetrics.metrics;
  }
  

export default async function MetricsPage() {
    //Get the rule findings metrics from the API
    const findingsMetrics: FindingsMetrics[] = await GetRuleFindingsMetrics();
	//Get the rule ids metrics from the API
	const ruleIdMetrics: FindingsMetrics[] = await GetRuleIdsMetrics();

    return (
        <div className="flex flex-row gap-4 flex-wrap">
            <div className="w-1/5 min-w-fit grow">
               <FindingsRadarChart title="# classifications" labels={findingsMetrics.map(({classification}) => classification)} values={findingsMetrics.map(({count}) => count)}/>
            </div>
			<div className="w-1/5 min-w-fit grow">
                <CustomBarChart labels={ruleIdMetrics.map(({classification}) => classification)} values={ruleIdMetrics.map(({count}) => count)} title="No. Findings per Rule"/>
            </div>
            {/* <div className="w-1/5 min-w-fit grow">
                <CustomPieChart labels={['Classified', 'Unclassified']} values={[classificationMetrics.classifiedCount, classificationMetrics.unclassifiedCount]} title="Classification Counts"/>
            </div> */}
        </div>
    );
}