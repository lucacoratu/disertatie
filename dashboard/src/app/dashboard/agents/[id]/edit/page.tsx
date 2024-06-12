import { constants } from "@/app/constants";
import AgentEditForm from "@/components/AgentEditForm";
import {
    HoverCard,
    HoverCardContent,
    HoverCardTrigger,
} from "@/components/ui/hover-card";
import { cookies } from "next/headers";

//Function to get the details of an agent
async function GetAgent(id: string): Promise<Agent|undefined> {
    const cookie = cookies().get('session');
    const URL = constants.apiBaseURL + "/agents/" + id;
    const res = await fetch(URL, {headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
    if (!res) {
        throw new Error("cannot get agent details");
    }
    try {
        const agent: Agent = await res.json();
        return agent;
    } catch(exception) {
        return undefined;
    } 
}

async function getMachines() : Promise<Machine[]> {
    const cookie = cookies().get('session');
    //Create the URL for fetching the machines
    const url = constants.apiBaseURL + "/machines" 
    //Revalidate the data once every 10 mins
    const res = await fetch(url, { next: { revalidate: 600 }, headers: {Cookie: `${cookie?.name}=${cookie?.value}`} });
    //Check if there was an error
    if(!res.ok) {
        throw new Error("could not get machines data");
    }
    
    const machines: MachineResponse = await res.json();
    //Return the data
    return machines.machines;
}

export default async function AgentEditPage({ params }: { params: { id: string } }) {
    //Get the agent details
    const agentId: string = params.id; 
    const agent = await GetAgent(agentId);
    const machines = await getMachines();

    return (
        <div className="p-6 w-2/3 h-[700px] bg-card min-w-[400px] m-auto rounded-xl border-2">
            <div className="flex flex-row justify-between items-center mb-[10px]">
                <h1 className="text-xl">Edit Agent</h1>
                <HoverCard>
                    <HoverCardTrigger>
                        <div className="dark:bg-darksurface-300 rounded-full w-[25px] text-center">
                            ?
                        </div>
                    </HoverCardTrigger>
                    <HoverCardContent className="flex flex-col items-start w-[400px]">
                        <h1 className="text-base">Editing an agent</h1>
                        <div className="flex flex-col text-xs">
                            <p>Here you can edit details about the agent.</p>
                            <p>Modify only the fields that are outdated, and submit the form.</p>
                        </div>
                    </HoverCardContent>
                </HoverCard>
            </div>
            <AgentEditForm agent={agent} machines={machines}/>
        </div>
    )
}