import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";

import {
    ContextMenu,
    ContextMenuContent,
    ContextMenuItem,
    ContextMenuTrigger,
} from "@/components/ui/context-menu";

import OsCard from "@/components/OsCard";
import Link from 'next/link';
import { Separator } from "@/components/ui/separator"
import { FC } from "react";

const AgentCard: FC<AgentProps> = ({agent}): JSX.Element => {
    return (
        <ContextMenu>
            <ContextMenuTrigger className="min-w-96 w-1/4 rounded-lg hover:shadow-lg hover:dark:bg-darksurface-100/[.8] dark:bg-darksurface-100 dark:border-darksurface-400">
                <Card>
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
            </ContextMenuTrigger>
            <ContextMenuContent>
                <ContextMenuItem className="dark:hover:bg-darksurface-400">
                    <Link href={`/agents/${encodeURIComponent(agent.id)}`}>
                        Details
                    </Link>
                </ContextMenuItem>
                <ContextMenuItem className="dark:hover:bg-darksurface-400">
                    <Link href={`/agents/${encodeURIComponent(agent.id)}/edit`}>
                        Edit
                    </Link>
                </ContextMenuItem>
                <Separator orientation="horizontal"/>
                <ContextMenuItem className="mt-1 bg-red-600 dark:hover:bg-red-600/[0.8]">Delete</ContextMenuItem>
            </ContextMenuContent>
        </ContextMenu>
    );
}

export default AgentCard;