import AgentCard from "@/components/AgentCard";
import { Button } from "@/components/ui/button";
import { Play } from "lucide-react"
import {constants} from "@/app/constants"

async function getAgents() : Promise<Agent[]> {
    //Create the URL where the data will be fetched from
    const url = constants.apiBaseURL + "/agents";
    //Revalidate the data once every 10 mins
    const res = await fetch(url, { next: { revalidate: 600 } });
    //Check if there was an error
    if(!res.ok) {
        throw new Error("could not get machines data");
    }
    const agents: AgentResponse = await res.json();
    //Return the data
    return agents.agents;
}

export default async function AgentsPage() {
    //Load the data from the server
    const agents: Agent[] = await getAgents();

    return (
        <main className="h-full w-full flex flex-col py-2 gap-2 px-2 dark:bg-darksurface-200">
            <Button className="self-end h-8">
                <Play className="mr-2 h-4 w-4"/>Deploy Agent
            </Button>
            <div className="flex flex-row gap-2">
                {agents && agents.map((agent) => (
                    <AgentCard key={agent.id} agent={agent}/>
                ))}
            </div>
        </main>
    )
}