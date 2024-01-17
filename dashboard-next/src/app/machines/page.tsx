import { constants } from "@/app/constants";
import { Button } from "@/components/ui/button";
import MachineCard from "@/components/MachineCard";
import { Server } from "lucide-react"

async function getMachines() : Promise<Machine[]> {
    //Create the URL for fetching the machines
    const url = constants.apiBaseURL + "/machines" 
    //Revalidate the data once every 10 mins
    const res = await fetch(url, { next: { revalidate: 600 } });
    //Check if there was an error
    if(!res.ok) {
        throw new Error("could not get machines data");
    }
    
    const machines: MachineResponse = await res.json();
    //Return the data
    return machines.machines;
}

export default async function MachinesPage() {
    //Load the data from the server
    const machines: Machine[] = await getMachines();

    return (
        <main className="h-full w-full flex flex-col py-2 gap-2 px-2 dark:bg-darksurface-200">
            <Button className="self-end h-8">
                <Server className="mr-2 h-4 w-4"/>Register Machine
            </Button>
            <MachineCard machine={machines[0]}/>
        </main>
    )
}