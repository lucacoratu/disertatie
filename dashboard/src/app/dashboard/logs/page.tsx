import { Button } from "@/components/ui/button";
import Link from "next/link";

export default async function Logs({ params }: { params: { id: string } }) {
    return (
        <>
            <Button className="self-end h-8">
                <Link href={`/dashboard/logs/classified`} className="flex flex-row gap-1 items-center">
                    {/* <Play className="mr-2 h-4 w-4"/>Deploy Agent */}
                    Classified Logs
                </Link>
            </Button>
            
        </>
    )
}