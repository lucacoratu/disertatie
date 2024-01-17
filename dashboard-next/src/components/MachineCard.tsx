import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
  } from "@/components/ui/card"

import { FC } from "react";

import OsCard from "@/components/OsCard"

const AgentCard: FC<MachineProps> = ({machine}): JSX.Element => {
    return (
        <Card className="w-1/4 rounded-lg hover:shadow-lg hover:dark:bg-darksurface-100/[.8] dark:bg-darksurface-100 dark:border-darksurface-400">
            <CardHeader>
                <CardTitle>{machine?.hostname}</CardTitle>
                <CardDescription className="flex flex-col gap-0">
                    <span>{machine?.id}</span>
                    <span>Machine holds {machine?.numberAgents} {machine?.numberAgents === 1 ? "agent": "agents"}</span>
                </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-row gap-4 justify-between items-center">
                <div>
                    <p>Details 1</p>
                    <p>Details 2</p>
                </div>
                <OsCard os={machine.os}/>
            </CardContent>
            {/* <CardFooter>
                <p>Card Footer</p>
            </CardFooter> */}
        </Card>
    );
}

export default AgentCard;