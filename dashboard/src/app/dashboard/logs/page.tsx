import { redirect } from "next/navigation";

export default async function Logs() {
    //Redirect to /logs/all
    redirect("/dashboard/logs/all");
    return (
        <>
        </>
    )
}