import { constants } from "@/app/constants";
import ExportButton from "@/components/ExportButton";
import { cookies } from "next/headers";

async function ExportLogs(agentId: string, format: string) {
	"use server";
	const cookie = cookies().get('session');
	const URL = `${constants.apiBaseURL}/agents/${agentId}/export-logs?format=${format}`;
	const res = await fetch(URL, {cache:"no-store", headers: {Cookie: `${cookie?.name}=${cookie?.value}`}});
}

export default async function AgentDetails({ params }: { params: { id: string } }) {
	const agentId = params.id;
	return (
		<>
			<ExportButton clickHandler={ExportLogs} agentId={agentId}></ExportButton>
		</>
	);
}