import { constants } from "@/app/constants";

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

	return res;
}
