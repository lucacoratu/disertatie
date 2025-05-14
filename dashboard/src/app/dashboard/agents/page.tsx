import AgentCard from "@/components/AgentCard";
import { Button } from "@/components/ui/button";
import { Play } from "lucide-react";
import {constants} from "@/app/constants";
import Link from "next/link";
import { cookies } from "next/headers";

async function getAgents() : Promise<Agent[]> {
    const cookie = cookies().get('session');
    //Create the URL where the data will be fetched from
    const url = constants.apiBaseURL + "/agents";
    //Revalidate the data once every 10 mins
    const res = await fetch(url, { cache: "no-store", headers: {Cookie: `${cookie?.name}=${cookie?.value}`} });
    //Check if there was an error
    if(!res.ok) {
        throw new Error("could not get agents data");
    }
    const agents: AgentResponse = await res.json();
    //Return the data
    return agents.agents;
}

async function GetConnectedAgents() : Promise<string[]> {
    const cookie = cookies().get('session');
    //Create the URL where the data will be fetched froms
    const url = constants.apiBaseURL + "/agents/connected";
    const res = await fetch(url, {cache: "no-store", headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
    if(!res.ok) {
        throw new Error("could not get connected agents");
    }
    const connAgents: ConnectedAgentsResponse = await res.json();
    return connAgents.agents;
}

export default async function AgentsPage() {
    //Load the data from the server
    const agents: Agent[] = await getAgents();
    //Get the connected agents
    const connAgents: string[] = await GetConnectedAgents();

    return (
        <>
            <Button className="self-end h-8">
                <Link href={`/dashboard/agents/deploy`} className="flex flex-row gap-1 items-center">
                    <Play className="mr-2 h-4 w-4"/>Deploy Agent
                </Link>
            </Button>
            <div className="flex flex-row gap-2">
                {
                    agents && agents.map(
                        (agent) => {
                            //Get the agent status from the connAgents
                            var isConn: boolean = connAgents?.indexOf(agent.id) === -1 ? false : true;
                            if(connAgents === null)
                                isConn = false;

                            return <AgentCard key={agent.id} agent={agent} agentStatus={isConn === true ? "online" : "offline"}/>
                        }
                    )
                }
            </div>
        </>
    )
}