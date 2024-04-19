"use client";

import Link from "next/link"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { useState } from "react"
import { z } from "zod"
import { constants } from "@/app/constants";
import { toast } from "sonner";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Loader2 } from "lucide-react"

const loginFormSchema = z.object({
  username: z.string().min(1, {
    message: "Username should be at least 1 character long",
  }),
  password: z.string().min(1, {
    message: "The password should be at least 1 character long",
  }),
})

export default function LoginForm() {
  const form = useForm<z.infer<typeof loginFormSchema>>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
        username: "",
        password: "",
      },
  })

  const [loading, setLoading] = useState(false);
  const [failed, setFailed] = useState(false);
  const [failMessage, setFailMessage] = useState("");

  async function onSubmit(values: z.infer<typeof loginFormSchema>) {
    setLoading(true);
    //Send the values to the server
    const requestBody = JSON.stringify(values);
    const res = await fetch(`/api/login`, {
        method: "POST",
        body: requestBody,
        headers: {
            "Content-Type": "application/json"
        },
    });

    if(!res.ok) {
      const resp: APIError = await res.json();
      
      //Display the error
      setFailed(true);
      setFailMessage(resp.message);
      setLoading(false);
      return
    }

    if (res.ok) {
      const resp: LoginResponse = await res.json();
      
      toast("Logged in", {
          description: `You will be redirected to dashboard shortly`,
      });

      //Get the cookie
      const token: string = resp.token;
      //Save the token in local storage
      localStorage.setItem("token", token);
    }

    setFailed(false);
    setLoading(false);
  }

  return (
    <main className="w-screen h-screen flex flex-1 flex-col md:gap-8 bg-muted">
    <Card className="mx-auto max-w-sb my-auto">
      <CardHeader>
        <CardTitle className="text-2xl text-center">Login</CardTitle>
        <CardDescription className="text-center">
          Enter your username below to login to your account
        </CardDescription>
        {failed &&
          <CardDescription>
            <div className="bg-destructive text-center py-3 text-foreground">
              {failMessage}
            </div>
          </CardDescription>
        }
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="grid gap-4">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                  <FormItem>
                      <FormLabel>Username</FormLabel>
                      <FormControl>
                          <Input className="dark:border-darksurface-400" placeholder="Username" {...field} />
                      </FormControl>
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
                          <Input type="password" className="dark:border-darksurface-400" placeholder="password" {...field} />
                      </FormControl>
                      <FormMessage />
                  </FormItem>
              )}
            />
            { !loading && <Button className="w-1/2 m-auto mt-[20px]" type="submit">Authenticate</Button> }
            { loading &&  
                <Button className="w-1/2 m-auto mt-[20px]" disabled>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Authenticating...
                </Button>
            }
          </form>
        </Form>
        {/* <div className="grid gap-4">
          <div className="grid gap-2">
            <Label htmlFor="username">Username</Label>
            <Input
              id="username"
              type="text"
              placeholder="michael"
              required
            />
          </div>
          <div className="grid gap-2">
            <div className="flex items-center">
              <Label htmlFor="password">Password</Label>
              <Link href="#" className="ml-auto inline-block text-sm underline">
                Forgot your password?
              </Link>
            </div>
            <Input id="password" type="password" placeholder="***********" required />
          </div>
          <Button type="submit" className="mt-4 w-full">
            Login
          </Button>
        </div> */}
        <div className="mt-4 text-center text-sm">
          Don&apos;t have an account?{" "}
          <Link href="#" className="underline">
            Sign up
          </Link>
        </div>
      </CardContent>
    </Card>
    </main>
  );
}