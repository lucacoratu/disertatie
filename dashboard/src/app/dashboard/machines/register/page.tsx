"use client"
 
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { useState } from "react"
import { z } from "zod"
import { toast } from "sonner"
import { Loader2 } from "lucide-react"
 
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
    HoverCard,
    HoverCardContent,
    HoverCardTrigger,
  } from "@/components/ui/hover-card";

import { useCookies } from 'next-client-cookies';
  

const osRegex = new RegExp('^(linux|windows)$');

const formSchema = z.object({
  hostname: z.string().min(2, {
    message: "Hostname must be at least 2 characters.",
  }),
  os: z.string().regex(osRegex, {
    message: "OS can only be windows or linux"
  }),
  ipAddress: z.string().ip({
    version: "v4",
    message: "Invalid IPv4 address",
  }),
  username: z.string().min(1, {
    message: "Username should be at least 1 character long",
  }),
  password: z.string().min(1, {
    message: "The password should be at least 1 character long",
  }),
  privateKey: z.string().startsWith("-----BEGIN", {
    message: "Private key should be in PEM format (generate using ssh-keygen)",
  })
})

export default function MachinesRegisterPage() {
    // 1. Define the form form.
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            hostname: "",
            os: "",
            ipAddress: "",
            username: "",
            password: "",
            privateKey: "",
        },
    })
    
    const [loading, setLoading] = useState(false);
    const cookies = useCookies();

    // 2. Define a submit handler.
    async function onSubmit(values: z.infer<typeof formSchema>) {
        const cookie = cookies.get('session');
        setLoading(true);
        //Send the values to the server
        const requestBody = JSON.stringify(values);
        const res = await fetch("/api/machines", {
            method: "POST",
            body: requestBody,
            headers: {
                "Content-Type": "application/json",
                Cookie: `sesison=${cookie}`
            },
        });

        const resp: MachineRegisterResponse = await res.json(); 
        
        toast("Machine created", {
            description: `A new machine has been created, UUID: ${resp.uuid}`,
            action: {
                label: "Undo",
                onClick: () => {},
            }
        });
        setLoading(false);
    }

    return (
        <div className="p-6 w-1/2 h-[700px] min-w-[400px] bg-card m-auto rounded-xl border-2">
            <div className="flex flex-row justify-between items-center mb-[10px]">
                <h1 className="text-xl">Register Machine</h1>
                <HoverCard>
                    <HoverCardTrigger>
                        <div className="dark:bg-darksurface-300 rounded-full w-[25px] text-center">
                            ?
                        </div>
                    </HoverCardTrigger>
                    <HoverCardContent className="flex flex-col items-start w-[400px]">
                        <h1 className="text-base">Registering a machine</h1>
                        <div className="flex flex-col text-xs">
                            <p>Register a new machine in the system where agents will be deployed.</p>
                            <p>You need to complete both the details and the credentials.</p>
                        </div>
                    </HoverCardContent>
                </HoverCard>
            </div>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 flex flex-col">
                    <Tabs defaultValue="details" className="w-full">
                        <TabsList className="dark:bg-darksurface-200 w-full">
                            <TabsTrigger className="dark:bg-darksurface-200 w-1/2" value="details">Details</TabsTrigger>
                            <TabsTrigger className="dark:bg-darksurface-200 w-1/2" value="credentials">Credentials</TabsTrigger>
                        </TabsList>
                        <TabsContent value="details">
                            <FormField
                            control={form.control}
                            name="hostname"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Hostname</FormLabel>
                                    <FormControl>
                                        <Input className="dark:border-darksurface-400" placeholder="ubuntu" {...field} />
                                    </FormControl>
                                    <FormDescription>
                                        This is the hostname of the machine
                                    </FormDescription>
                                    <FormMessage  />
                                </FormItem>
                            )}
                            />
                            <FormField
                                control={form.control}
                                name="os"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Operating System</FormLabel>
                                        <FormControl>
                                            <Input className="dark:border-darksurface-400" placeholder="windows, linux" {...field} />
                                        </FormControl>
                                        <FormDescription>
                                            The OS of the machine
                                        </FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="ipAddress"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>IP Address</FormLabel>
                                        <FormControl>
                                            <Input className="dark:border-darksurface-400" placeholder="192.168.0.2" {...field} />
                                        </FormControl>
                                        <FormDescription>
                                            The IPv4 address of the machine
                                        </FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </TabsContent>
                        <TabsContent value="credentials">
                            <FormField
                                control={form.control}
                                name="username"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Username</FormLabel>
                                        <FormControl>
                                            <Input className="dark:border-darksurface-400" placeholder="username" {...field} />
                                        </FormControl>
                                        <FormDescription>
                                            The username used to connect to SSH
                                        </FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="password"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Password</FormLabel>
                                        <FormControl>
                                            <Input type="password" className="dark:border-darksurface-400" placeholder="SSH password" {...field} />
                                        </FormControl>
                                        <FormDescription>
                                            The password used to connect to SSH
                                        </FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="privateKey"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Private key</FormLabel>
                                        <FormControl>
                                            <Input className="dark:border-darksurface-400" placeholder="SSH private key" {...field} />
                                        </FormControl>
                                        <FormDescription>
                                            The private key used to connect to SSH (PEM format)
                                        </FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            { !loading && <Button className="w-1/2 ml-auto mt-[20px]" type="submit">Submit</Button> }
                            { loading &&  
                                <Button className="w-1/2 ml-auto mt-[20px]" disabled>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Creating machine
                                </Button>
                            }
                        </TabsContent>
                    </Tabs>
                </form>
            </Form>
      </div>
    );
}