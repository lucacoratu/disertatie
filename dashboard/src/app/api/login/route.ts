import { constants } from "@/app/constants";
import { cookies } from 'next/headers';

export async function POST(request: Request) {
	//Send the data to the API
	const requestBody: string = await request.text();
	const URL: string = constants.apiBaseURL + "/auth/login";
	const res = await fetch(URL, {
		method: "POST",
		body: requestBody,
		headers: {
			"Content-Type": "application/json"
		},
	});
	const setcookies: string[] = res.headers.getSetCookie();
	console.log(setcookies);
	cookies().set('token', setcookies[0]);

	return res;
}
