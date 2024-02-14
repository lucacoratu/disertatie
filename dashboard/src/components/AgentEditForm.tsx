"use client"

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { z } from "zod";
import { toast } from "sonner";
import { Loader2, ArrowBigRight } from "lucide-react";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { DisplayAPIErrorToast } from "@/lib/utils";

//The form schema
const formSchema = z.object({
    name: z.string().min(2, {
      message: "Name must be at least 2 characters.",
    }),
    machine: z.string({
        required_error: "Please select the machine the agent is deployed on.",
    }),
    listeningProtocol: z.string({
        required_error: "Please select the protocol the agent listens on",
    }),
    ipAddress: z.string().ip({
        version: "v4",
        message: "IP address is required",
    }),
    port: z.number().min(1, {
        message:"Port should be >= 1"
    }).max(65535, {
        message: "Port should be <= 65535"
    }),
    webServerProtocol: z.string({
        required_error: "Please select the protocol the webserver listens on",
    }),
    webServerIpAddress: z.string().ip({
        version: "v4",
        message: "IP address is required",
    }),
    webServerPort: z.number().min(1, {
        message:"Port should be >= 1"
    }).max(65535, {
        message: "Port should be <= 65535"
    })
});

type AgentEditFormProps = {
    agent: Agent | undefined;
    machines: Machine[] | undefined;
}



export default function AgentEditForm(props: AgentEditFormProps) {
    const agent = props.agent;
    const machines = props.machines; 

    // 1. Define the form form.
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            name: agent?.name || "",
            machine: machines?.filter((machine) => machine.id == agent?.machineId)[0].id,
            ipAddress: agent?.listeningAddress,
            listeningProtocol: agent?.listeningProtocol.toUpperCase(),
            port: Number(agent?.listeningPort),
            webServerProtocol: agent?.forwardServerProtocol.toUpperCase(),
            webServerIpAddress: agent?.forwardServerAddress,
            webServerPort: Number(agent?.forwardServerPort),
        },
    });

    const [loading, setLoading] = useState(false);

    // Form submit handler
    async function onSubmit(values: z.infer<typeof formSchema>) {
        setLoading(true);
        //Create the structure which will be sent to the server
        const updateAgent: UpdateAgent = {
            name: values.name,
            listeningProtocol: values.listeningProtocol.toLowerCase(),
            listeningAddress: values.ipAddress,
            listeningPort: values.port,
            forwardServerProtocol: values.webServerProtocol.toLowerCase(),
            forwardServerAddress: values.webServerIpAddress,
            forwardServerPort: values.webServerPort,
            machineId: values.machine,
        }
        const requestBody = JSON.stringify(updateAgent);

        console.log(requestBody);

        const res = await fetch(`/api/agents/${agent?.id}`, {
            method: "PUT",
            body: requestBody,
            headers: {
                "Content-Type": "application/json"
            },
        });

        if (!res.ok) {
            //Get the error from the server
            try {
                const apiErr: APIError = await res.json();
                DisplayAPIErrorToast("Failed to modify", `Agent ${agent?.id} could not be modified`, apiErr);
                setLoading(false);
                return;
            } catch(exception: any) {
                DisplayAPIErrorToast("Failed to modify", `Agent ${agent?.id} could not be modified`, exception);
                setLoading(false);
                return;
            }
        }

        toast("Agent updated",{
            description: `Updated agent ${agent?.id}`,
            action: {
                label: "Ok",
                onClick: () => {},
            },
        });

        setLoading(false);
    }

    return (
        <>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-2 flex flex-col">
                    <FormField
                        control={form.control}
                        name="name"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Name</FormLabel>
                                <FormControl>
                                    <Input className="dark:border-darksurface-400" placeholder="Test" {...field} />
                                </FormControl>
                                <FormDescription>
                                    The name of the agent
                                </FormDescription>
                                <FormMessage  />
                            </FormItem>
                        )}
                    />
                     <FormField
                        control={form.control}
                        name="machine"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Machine</FormLabel>
                                <Select  onValueChange={field.onChange} defaultValue={field.value}>
                                    <FormControl>
                                        <SelectTrigger className="dark:border-darksurface-400">
                                            <SelectValue placeholder="Select the machine where agent is deployed" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <SelectContent>
                                        {machines?.map((machine) => {
                                            return <>
                                                <SelectItem className="dark:hover:bg-darksurface-400 dark:selected:bg-darksurface-400" key={machine.id} value={machine.id}>{machine.hostname} - {machine.id}</SelectItem>
                                            </>;
                                        })}
                                    </SelectContent>
                                </Select>
                                <FormDescription>
                                    Machine the agent is deployed on
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                        />
                        <div className="flex flex-col gap-2">
                            <div className="flex flex-row items-center justify-around">
                                <h1>Agent config</h1>
                                <h1>WebServer config</h1>
                            </div>
                            <div className="flex flex-row gap-0 items-center justify-around">
                                <div className="flex flex-row gap-1 items-center justify-center">
                                    <FormField
                                        control={form.control}
                                        name="listeningProtocol"
                                        render={({ field }) => (
                                            <FormItem>
                                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                                    <FormControl>
                                                        <SelectTrigger className="dark:border-darksurface-400 grow">
                                                            <SelectValue placeholder="Select the machine where agent is deployed" />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        <SelectItem className="dark:hover:bg-darksurface-400 dark:selected:bg-darksurface-400" value="HTTP">HTTP</SelectItem>
                                                        <SelectItem className="dark:hover:bg-darksurface-400 dark:selected:bg-darksurface-400" value="HTTPS">HTTPS</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="ipAddress"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormControl>
                                                    <Input className="dark:border-darksurface-400 grow" placeholder="0.0.0.0" {...field} />
                                                </FormControl>
                                                <FormMessage  />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="port"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormControl>
                                                    <Input className="dark:border-darksurface-400 grow" placeholder="8080" {...field} />
                                                </FormControl>
                                                <FormMessage  />
                                            </FormItem>
                                        )}
                                    />
                                </div>
                                <ArrowBigRight className="w-8 h-8" />
                                <div className="flex flex-row gap-1 items-center justify-center">
                                    <FormField
                                        control={form.control}
                                        name="webServerProtocol"
                                        render={({ field }) => (
                                            <FormItem>
                                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                                    <FormControl>
                                                        <SelectTrigger className="dark:border-darksurface-400 grow">
                                                            <SelectValue placeholder="Select the machine where agent is deployed" />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        <SelectItem className="dark:hover:bg-darksurface-400 dark:selected:bg-darksurface-400" value="HTTP">HTTP</SelectItem>
                                                        <SelectItem className="dark:hover:bg-darksurface-400 dark:selected:bg-darksurface-400" value="HTTPS">HTTPS</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="webServerIpAddress"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormControl>
                                                    <Input className="dark:border-darksurface-400 grow" placeholder="0.0.0.0" {...field} />
                                                </FormControl>
                                                <FormMessage  />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="webServerPort"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormControl>
                                                    <Input className="dark:border-darksurface-400 grow" placeholder="8080" {...field} />
                                                </FormControl>
                                                <FormMessage  />
                                            </FormItem>
                                        )}
                                    />
                                </div>
                            </div>
                        </div>
                    { !loading && <Button className="w-1/2 m-auto mt-[30px]" type="submit">Submit</Button> }
                    { loading &&  
                        <Button className="w-1/2 m-auto mt-[30px]" disabled>
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Updating agent
                        </Button>
                    }
                </form>
            </Form>
        </>
    );
}