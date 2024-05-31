import { constants } from "@/app/constants";
import { NextApiRequest } from "next";
import { cookies } from 'next/headers'

export async function PUT(request: Request, context: any) {
    //Get the id from the context
    const {id} = context.params;
    //Get the cookie from the body
    const jwtToken = cookies().get("session")?.value;
    const requestBody: string = await request.text();
    //Send the data to the API
    const URL = constants.apiBaseURL + `/agents/${id}`;
    const res = await fetch(URL, {
        method: "PUT",
        body: requestBody,
        headers: {
            "Content-Type": "application/json",
            "Cookie": `session=${jwtToken}`,
        },
    });
  	return res;
} 
