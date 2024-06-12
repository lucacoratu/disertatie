"use client"

import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";

import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    DialogFooter,
    DialogClose,
} from "@/components/ui/dialog"

import {
    ContextMenu,
    ContextMenuContent,
    ContextMenuItem,
    ContextMenuTrigger,
} from "@/components/ui/context-menu";

import { FC } from "react";
import { Separator } from "@/components/ui/separator"
import {Button} from "@/components/ui/button";
import OsCard from "@/components/OsCard";
import Link from 'next/link';
import {Trash, Pencil, GanttChartSquare } from "lucide-react";
import { toast } from "sonner";

async function DeleteClicked(machine: Machine) {
    //Send the request to the backend
    const res = await fetch(`/api/machines/${machine.id}`, {
        method: "DELETE"
    });

    if (!res.ok) {
        //Get the error from the server
        try {
            const apiErr: APIError = await res.json();

            toast.error("Failed to delete", {
                description: <div className="flex flex-col gap-1"><p>Machine {machine.hostname} could not be deleted</p><p>Code: {apiErr.code}, Message: {apiErr.message}</p></div>,
                classNames: {
                    error: 'bg-red-400',
                }
            });
    
            return;
        } catch(exception) {
            toast.error("Failed to delete", {
                description: <div className="flex flex-col gap-1"><p>Machine {machine.hostname} could not be deleted</p><p>Unknown error</p></div>,
                classNames: {
                    error: 'bg-red-400',
                }
            });
        }
    }

    toast("Machine deleted",{
        description: `Deleted machine ${machine.id}`,
        action: {
            label: "Ok",
            onClick: () => {},
        },
    });

    //Refresh the page
}

const MachineCard: FC<MachineProps> = ({machine}): JSX.Element => {
    return (
        <Dialog>
            <ContextMenu>
                <ContextMenuTrigger className="min-w-96 w-1/4 rounded-lg hover:shadow-lg hover:dark:bg-darksurface-100/[.8] dark:bg-darksurface-100 dark:border-darksurface-400">
                    <Link href={`/machines/${encodeURIComponent(machine?.id)}`}>
                        <Card>
                            <CardHeader>
                                <CardTitle>{machine?.hostname}</CardTitle>
                                <CardDescription className="flex flex-col gap-0">
                                    <span>{machine?.id}</span>
                                    <span>Machine holds {machine?.numberAgents} {machine?.numberAgents === 1 ? "agent": "agents"}</span>
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="flex flex-row gap-4 justify-between items-center">
                                <div>
                                    <p>{machine?.ipAddresses.length} network interfaces</p>
                                    <p>Details 2</p>
                                </div>
                                <OsCard os={machine.os}/>
                            </CardContent>
                            {/* <CardFooter>
                                <p>Card Footer</p>
                            </CardFooter> */}
                        </Card>
                    </Link>
                </ContextMenuTrigger>
                <ContextMenuContent>
                    <ContextMenuItem className="dark:hover:bg-darksurface-400 flex flex-row items-center gap-4">
                        <GanttChartSquare className="w-4 h-4"/>Details
                    </ContextMenuItem>
                    <ContextMenuItem className="dark:hover:bg-darksurface-400 flex flex-row items-center gap-4">
                        <Pencil className="w-4 h-4"/>Edit
                    </ContextMenuItem>
                    <Separator orientation="horizontal"/>
                        <ContextMenuItem className="mt-1 bg-red-600 dark:hover:bg-red-600/[0.8]">
                            <DialogTrigger className="flex flex-row items-center gap-4"><Trash className="w-4 h-4"/> Delete</DialogTrigger>
                        </ContextMenuItem>
                </ContextMenuContent>
            </ContextMenu>
            <DialogContent>
                <DialogHeader>
                <DialogTitle>Are you absolutely sure?</DialogTitle>
                <DialogDescription className="flex flex-col gap-0">
                    <p>This action will permanently delete the selected machine.</p>
                    <p>By clicking delete you will lose all data associated with this machine</p>
                </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <DialogClose asChild>
                        <Button variant="outline">Cancel</Button>
                    </DialogClose>
                    <DialogClose>
                        <Button className="bg-red-600 dark:hover:bg-red-600/[0.8]" onClick={() => DeleteClicked(machine)}>Delete</Button>
                    </DialogClose>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

export default MachineCard;