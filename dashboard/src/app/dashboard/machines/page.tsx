import { constants } from "@/app/constants";
import { Button } from "@/components/ui/button";
import MachineCard from "@/components/MachineCard";
import { Server } from "lucide-react";
import Link from "next/link";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { cookies } from "next/headers";

async function getMachines() : Promise<Machine[]> {
    const cookie = cookies().get('session');
    //Create the URL for fetching the machines
    const url = constants.apiBaseURL + "/machines" 
    //Revalidate the data once every 10 mins
    const res = await fetch(url, { cache: "no-store", headers: {Cookie: `${cookie?.name}=${cookie?.value}`} });
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
        <>
            <Button className="self-end h-8">
                <Link href={`/dashboard/machines/register`} className="flex flex-row gap-1 items-center">
                    <Server className="mr-2 h-4 w-4"/>Register Machine
                </Link>
            </Button>
            <ScrollArea>
                <div className="flex flex-row gap-2">
                    {machines && machines.map((machine) => {
                        return <MachineCard key={machine.id} machine={machine}/>
                    })}
                </div>
                <ScrollBar className="mt-[5px]" orientation="horizontal" />
            </ScrollArea>
        </>
    )
}