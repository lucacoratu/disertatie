import { constants } from "@/app/constants";

export async function DELETE(request: Request, context: any) {
	//Send the delete request to the API
    const { id } = context.params;
    const URL = `${constants.apiBaseURL}/machines/${id}`;
    const res = await fetch(URL, {
        method: "DELETE",
    });

    //Send the response to the client
    return res;
}