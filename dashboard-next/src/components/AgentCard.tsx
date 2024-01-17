import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
  } from "@/components/ui/card"

import OsCard from "@/components/OsCard"
import Link from 'next/link'

import { FC } from "react";

const AgentCard: FC<AgentProps> = ({agent}): JSX.Element => {
    return (
        <Card className="min-w-96 w-1/4 rounded-lg hover:shadow-lg hover:dark:bg-darksurface-100/[.8] dark:bg-darksurface-100 dark:border-darksurface-400">
            <Link href={`/agents/${encodeURIComponent(agent.id)}`}>
                <CardHeader>
                    <CardTitle>{agent?.name ? agent?.name : "No name"}</CardTitle>
                    <CardDescription className="flex flex-col gap-0">
                        <span>{agent?.id}</span>
                        <span>Deployed on {agent?.machineHostname}</span>
                    </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-row gap-4 justify-between items-center">
                    <div>
                        <p>{agent?.machineOs}</p>
                        <p>{agent?.listeningProtocol}://{agent?.listeningAddress}:{agent?.listeningPort} &#8594; {agent?.forwardServerProtocol}://{agent?.forwardServerAddress}:{agent?.forwardServerPort}</p>
                    </div>
                    <OsCard os={agent?.machineOs}/>
                </CardContent>
                {/* <CardFooter>
                    <p>Card Footer</p>
                </CardFooter> */}
            </Link>
        </Card>
    );
}

export default AgentCard;