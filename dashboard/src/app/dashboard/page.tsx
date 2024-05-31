
import Link from "next/link";
import {
  Ghost,
  Server,
  ArrowUpRight,
  ScrollText,
  Radar,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

import {constants} from "@/app/constants";
import RuleFindingPreview from "@/components/ui/rulefinding-preview";
import { cookies } from "next/headers";


async function GetRecentLogs() {
  const cookie = cookies().get('session');
  //Create the URL where the logs will be fetched from
  const URL = `${constants.apiBaseURL}/logs/recent`;
  //Fetch the data (revalidate after 10 minutes)
  const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
  //Check if an error occured
  if(!res.ok) {
    throw new Error("could not load logs");
  }
  //Parse the json data
  const logsResponse: LogsShortElasticResponse = await res.json();
  return logsResponse;
}

async function GetRecentClassifiedLogs() {
  const cookie = cookies().get('session');
  //Create the URL where the logs will be fetched from
  const URL = `${constants.apiBaseURL}/logs/recent-classified`;
  //Fetch the data (revalidate after 10 minutes)
  const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
  //Check if an error occured
  if(!res.ok) {
    throw new Error("could not load logs");
  }
  //Parse the json data
  const logsResponse: LogsShortElasticResponse = await res.json();
  return logsResponse;
}

async function GetMachinesStatistics() {
  const cookie = cookies().get('session');
  //Create the URL where the logs will be fetched from
  const URL = `${constants.apiBaseURL}/machines/metrics`;
  //Fetch the data (revalidate after 10 minutes)
  const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
  //Check if an error occured
  if(!res.ok) {
    throw new Error("could not load machines metrics");
  }
  //Parse the json data
  const machinesStatisticsResponse: MachinesStatisticsResponse = await res.json();
  return machinesStatisticsResponse;
}

async function GetTotalLogsCount() {
  const cookie = cookies().get('session');
  //Create the URL where the logs will be fetched from
  const URL = `${constants.apiBaseURL}/logs/count`;
  //Fetch the data (revalidate after 10 minutes)
  const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
  //Check if an error occured
  if(!res.ok) {
    throw new Error("could not load machines metrics");
  }
  //Parse the json data
  const logCount: LogCountResponse = await res.json();
  return logCount.count;
}

async function GetRuleFindingsMetrics() {
  const cookie = cookies().get('session');
    //Create the URL where the logs will be fetched from
    const URL = `${constants.apiBaseURL}/findings/rule/metrics`;
    //Fetch the data (revalidate after 10 minutes)
    const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
    //Check if an error occured
    if(!res.ok) {
      throw new Error("could not load rule findings metrics");
    }
    //Parse the json data
    const findingsMetrics: FindingsMetricsResponse = await res.json();
    return findingsMetrics.metrics;
}

async function GetRuleIdsMetrics() {
  const cookie = cookies().get('session');
  //Create the URL where the logs will be fetched from
  const URL = `${constants.apiBaseURL}/findings/rule/id-metrics`;
  //Fetch the data (revalidate after 10 minutes)
  const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
  //Check if an error occured
  if(!res.ok) {
    throw new Error("could not load rule ids metrics");
  }
  //Parse the json data
  const findingsMetrics: FindingsMetricsResponse = await res.json();
  return findingsMetrics.metrics;
}

async function GetCountAgents() {
  const cookie = cookies().get('session');
  //Create the URL the logs will be fetched from
  const URL = `${constants.apiBaseURL}/agents/count`;
  const res = await fetch(URL, {next: {revalidate: 600}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
  if (!res.ok) {
    throw new Error("could not load agents count");
  }
  const countResponse: AgentCountResponse = await res.json();
  return countResponse.count;
}

async function GetFindingsCountMetrics() {
  const cookie = cookies().get('session');
  //Create the URL the metrics will be pulled from
  const URL = `${constants.apiBaseURL}/findings/count-metrics`;
  const res = await fetch(URL, {next: {revalidate: 0}, headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
  if(!res.ok) {
    throw new Error("could not load findings count metrics");
  }
  const findingsCountMetricsResponse: FindingsCountMetricsResponse = await res.json();
  return findingsCountMetricsResponse.metrics;
}

export default async function DashboardHome() {
  //Get the logs from the api
  const logs: LogsShortElasticResponse = await GetRecentLogs();
  const classifiedLogs: LogsShortElasticResponse = await GetRecentClassifiedLogs();
  // console.log(classifiedLogs.logs[0].ruleFindings);
  const machineStatistics: MachinesStatisticsResponse = await GetMachinesStatistics();
  const totalCountLogs: number = await GetTotalLogsCount();
  const ruleFindingsMetrics: FindingsMetrics[] = await GetRuleFindingsMetrics();
  const ruleIdsMetrics: FindingsMetrics[] = await GetRuleIdsMetrics();
  const countAgents: number = await GetCountAgents();
  const findingsCountMetrics: FindingsCountMetrics = await GetFindingsCountMetrics();

  //Get the labels and the data from ruleFindingsMetrics
  const ruleFindingsMetricsLabels: string[] = ruleFindingsMetrics.map((metric) => metric.classification);
  const ruleFindingsMetricsCounts: number[] = ruleFindingsMetrics.map((metric) => metric.count);

  return (
    <main className="flex flex-1 flex-col gap-4 p-4 md:gap-8 md:p-6">
        <div className="grid gap-4 md:grid-cols-2 md:gap-8 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Machines
              </CardTitle>
              <Server className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {machineStatistics.totalMachines}
              </div>
              <p className="text-xs text-muted-foreground">
                {machineStatistics.totalInterfaces} network interfaces
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Agents
              </CardTitle>
              <Ghost className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{countAgents}</div>
              <p className="text-xs text-muted-foreground">
                1000 total commands ran
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Logs</CardTitle>
              <ScrollText className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{totalCountLogs}</div>
              <p className="text-xs text-muted-foreground">
                +19% from last month
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Findings</CardTitle>
              <Radar className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent className="flex justify-between">
              <div>
                <div className="text-2xl font-bold text-left">{findingsCountMetrics.ruleFindingsCount}</div>
                <p className="text-xs text-muted-foreground">
                  total rule findings
                </p>
              </div>
              <div>
                <div className="text-2xl font-bold text-right">{findingsCountMetrics.findingsCount}</div>
                <p className="text-xs text-muted-foreground">
                  total findings
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
        <div className="grid gap-4 md:gap-8 lg:grid-cols-2 xl:grid-cols-3">
          <Card className="xl:col-span-2">
            <CardHeader className="flex flex-row items-center">
              <div className="grid gap-2">
                <CardTitle>Recent Logs</CardTitle>
                <CardDescription>
                  Recent logs collected by agents.
                </CardDescription>
              </div>
              <Button asChild size="sm" className="ml-auto gap-1">
                <Link href="/dashboard/logs">
                  View All
                  <ArrowUpRight className="h-4 w-4" />
                </Link>
              </Button>
            </CardHeader>
            <CardContent>
              <Table className="text-sm">
                <TableHeader>
                  <TableRow>
                    <TableHead className="text-left max-w-30">
                      Agent
                    </TableHead>
                    <TableHead className="text-center">
                      Date
                    </TableHead>
                    <TableHead className="text-center max-w-14">
                      Method
                    </TableHead>
                    <TableHead className="text-center">
                      URL
                    </TableHead>
                    <TableHead className="text-right max-w-16">
                      Status Code
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {logs.logs.map((log) => {
                    //Convert the date from unix timestamp to locale date
                    const logDate: Date = new Date(log.timestamp * 1000);
                    const requestPreviewParts: string[] = log.request_preview.split(' '); 
                    const responsePreviewParts: string[] = log.response_preview.split(' '); 
                    return (
                      <TableRow key={log.id}>
                        <TableCell className="text-left max-w-28">
                          <div className="font-medium">{log.agentName}</div>
                          <div className="text-sm truncate text-ellipsis overflow-hidden text-muted-foreground md:inline">
                            {log.agentId}
                          </div>
                        </TableCell>
                        <TableCell className="text-center max-w-14">
                          {logDate.toLocaleString()}
                        </TableCell>
                        <TableCell className="text-center max-w-10">
                          {requestPreviewParts[0]}
                        </TableCell>
                        <TableCell className="text-left truncate">
                          {requestPreviewParts[1]}
                        </TableCell>
                        <TableCell className="text-right max-w-16">
                          {responsePreviewParts[1]}
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          <div className="grid grid-cols-1 gap-5 justify-stretch">
            <Card>
              <CardHeader>
                <CardTitle>Findings Statistics</CardTitle>
                <CardDescription className="flex flex-row justify-between">
                    <span> Rule Findings</span>
                    <span> Findings</span>
                  </CardDescription>
              </CardHeader>
              <CardContent className="flex flex-row gap-4">
                {/* <div className="w-64 h-64">
                  <FindingsRadarChart labels={ruleFindingsMetricsLabels} values={ruleFindingsMetricsCounts} title="# of rule findings" />
                </div> */}
                <div className="flex flex-col gap-3 w-1/3">
                  {ruleFindingsMetrics.map((metric, index) => {
                    return (
                      <div key={index} className="flex items-center gap-4">
                        <Badge className="rounded">{metric.classification}</Badge>
                        <div className="ml-auto font-medium">{metric.count}</div>
                      </div>
                    );
                  })}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Rule Statistics</CardTitle>
                <CardDescription>Number of matches for each rule</CardDescription>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="text-left">
                        Rule ID
                      </TableHead>
                      <TableHead className="text-right">
                        Detections
                      </TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {ruleIdsMetrics.map((metric, index) => {
                      return (
                        <TableRow key={index}>
                          <TableCell className="text-left">{metric.classification}</TableCell>
                          <TableCell className="text-right ml-auto font-medium">{metric.count}</TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </div>
        </div>

        <div className="grid gap-4 md:gap-8 lg:grid-cols-1 xl:grid-cols-1">
            <Card className="xl:col-span-2">
              <CardHeader className="flex flex-row items-center">
                <div className="grid gap-2">
                  <CardTitle>Recent Classified Logs</CardTitle>
                  <CardDescription>
                    Recent logs collected by agents which were classified using the rules.
                  </CardDescription>
                </div>
                <Button asChild size="sm" className="ml-auto gap-1">
                  <Link href="/dashboard/logs/classified">
                    View All
                    <ArrowUpRight className="h-4 w-4" />
                  </Link>
                </Button>
              </CardHeader>
              <CardContent>
                <Table className="text-sm">
                  <TableHeader>
                    <TableRow>
                      <TableHead className="text-left max-w-32">
                        Agent
                      </TableHead>
                      <TableHead className="text-center">
                        Date
                      </TableHead>
                      <TableHead className="text-center max-w-10">
                        Method
                      </TableHead>
                      <TableHead className="text-center max-w-64">
                        URL
                      </TableHead>
                      <TableHead className="text-center max-w-10">
                        Status Code
                      </TableHead>
                      <TableHead className="text-right max-w-32">
                        Classifications
                      </TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {classifiedLogs.logs.map((log) => {
                      //Convert the date from unix timestamp to locale date
                      const logDate: Date = new Date(log.timestamp * 1000);
                      const requestPreviewParts: string[] = log.request_preview.split(' '); 
                      const responsePreviewParts: string[] = log.response_preview.split(' ');
                      const findings: RuleFinding[] = log.ruleFindings;
                      return (
                        <TableRow key={log.id}>
                          <TableCell className="text-left max-w-32">
                            <div className="font-medium">{log.agentName}</div>
                            <div className="text-sm truncate text-ellipsis overflow-hidden text-muted-foreground md:inline">
                              {log.agentId}
                            </div>
                          </TableCell>
                          <TableCell className="text-center">
                            {logDate.toLocaleString()}
                          </TableCell>
                          <TableCell className="text-center max-w-10">
                            {requestPreviewParts[0]}
                          </TableCell>
                          <TableCell className="text-left truncate">
                            {requestPreviewParts[1]}
                          </TableCell>
                          <TableCell className="text-center max-w-32">
                            {responsePreviewParts[1]}
                          </TableCell>
                          <TableCell className="text-right flex flex-row gap-5 items-center">
                            {findings.map((finding) => {
                              return <RuleFindingPreview key={finding.request?.id} ruleFinding={finding}/>
                            })}
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </div>
      </main>
  )
}
