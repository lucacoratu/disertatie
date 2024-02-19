"use client";

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
import { Separator } from "@/components/ui/separator";
import { FC, useEffect } from "react";
import {Trash, Pencil, GanttChartSquare } from "lucide-react";
import WebSocketConnection from "@/types/websocket";
import {WSNotification} from "@/types/websocket_types";
import { useState } from "react";

const AgentCard: FC<AgentProps> = ({agent}): JSX.Element => {
    const [status, setStatus] = useState("offline");

    useEffect(() => {
        //Connect to the websocket of the API
        const wsConnection: WebSocketConnection = WebSocketConnection.getInstance();

        //Add connect callback
        wsConnection.addAgentConnectedCallback((message: string) => {
            //const notif: WSNotification = JSON.parse(message);
            setStatus("online");
        });
        
        //Add disconnect callback
        wsConnection.addAgentDisconnectedCallback((message: string) => {
            //const notif: WSNotification = JSON.parse(message);
            setStatus("offline");
        });
    }, [])



    return (
        <ContextMenu>
            <ContextMenuTrigger className="min-w-96 w-1/4 rounded-lg hover:shadow-lg hover:dark:bg-darksurface-100/[.8] dark:bg-darksurface-100 dark:border-darksurface-400">
                <Card>
                    <Link href={`/agents/${encodeURIComponent(agent.id)}`}>
                        <CardHeader>
                            <CardTitle>{agent?.name ? agent?.name : "No name"}</CardTitle>
                            <CardDescription className="flex flex-row items-center justify-between">
                                <span className="flex flex-col items-start">
                                    <span>{agent?.id}</span>
                                    <span>Deployed on {agent?.machineHostname}</span>
                                </span>
                                {
                                    status === "offline" && <span className="rounded-full w-5 h-5 mr-2 bg-red-500"/>
                                }
                                {
                                    status === "online" && <span className="rounded-full w-5 h-5 mr-2 bg-green-500"/>
                                }
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
                    <Link className="flex flex-row items-center gap-4" href={`/agents/${encodeURIComponent(agent.id)}`}>
                        <GanttChartSquare className="w-4 h-4"/>Details
                    </Link>
                </ContextMenuItem>
                <ContextMenuItem className="dark:hover:bg-darksurface-400 flex flex-row items-center gap-4">
                    <Link className="flex flex-row items-center gap-4" href={`/agents/${encodeURIComponent(agent.id)}/edit`}>
                        <Pencil className="w-4 h-4"/>Edit
                    </Link>
                </ContextMenuItem>
                <Separator orientation="horizontal"/>
                <ContextMenuItem className="mt-1 bg-red-600 dark:hover:bg-red-600/[0.8] flex flex-row items-center gap-4">
                    <Trash className="w-4 h-4"/> Delete
                </ContextMenuItem>
            </ContextMenuContent>
        </ContextMenu>
    );
}

export default AgentCard;