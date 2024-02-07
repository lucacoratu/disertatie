import { constants } from "@/app/constants";
import {
    HoverCard,
    HoverCardContent,
    HoverCardTrigger,
} from "@/components/ui/hover-card";

//Function to get the details of an agent
async function GetAgent(id: string): Promise<Agent|undefined> {
    const URL = constants.apiBaseURL + "/agents/" + id;
    const res = await fetch(URL);
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

export default async function AgentEditPage({ params }: { params: { id: string } }) {
    //Get the agent details
    const agentId: string = params.id; 
    const agent = await GetAgent(agentId);

    return (
        <div className="p-6 w-1/2 h-[700px] min-w-[400px] dark:bg-darksurface-100 m-auto rounded-xl dark:border-darksurface-400 border-2">
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
        </div>
    )
}