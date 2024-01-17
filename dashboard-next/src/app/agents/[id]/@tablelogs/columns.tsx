"use client"

import { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown, MoreHorizontal } from "lucide-react"
import FindingPreview from "@/components/ui/finding-preview"

import Link from "next/link";

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

//Create another type which will define the columns of the table
export type LogColumn = {
    id: string,
    agentId: string,
    remoteIp: string,
    requestMethod: string,
    url: string,
    response: string,
    timestamp: string,
    findings: Finding[],
    findingsClassificationString: FindingClassificationString[],
}

export const columns: ColumnDef<LogColumn>[] = [
    {
        accessorKey: "remoteIp",
        header: ({ column }) => {
            return (
              <Button
                variant="ghost"
                onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
              >
                Remote IP
                <ArrowUpDown className="ml-2 h-4 w-4" />
              </Button>
            )
          },
    },
    {
        accessorKey: "requestMethod",
        header: "Request Method"
    },
    {
        accessorKey: "url",
        header: "URL",
    },
    {
        accessorKey: "response",
        header: "Response",
    },
    {
        accessorKey: "timestamp",
        header: "Timestamp",
    },
    {
        accessorKey: "findings",
        header: "Findings",
        cell: ({row}) => {
            const log = row.original;
            //console.log(log);
            if(log.findings != null) {
              return (
                  <div className="overflow-hidden flex flex-row gap-2 max-w-36">
                    {log.findings.map((finding: Finding) => (
                        <FindingPreview key={finding.request.id} finding={finding} findingsClassificationString={log.findingsClassificationString}/>
                    ))
                    }
                  </div>
              );
            } else {
              return (
                <>
                </>
              );
            }
        }
    },
    {
        id: "actions",
        header: "Actions",
        cell: ({ row }) => {
          const log = row.original
     
          return (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                  <span className="sr-only">Open menu</span>
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="dark:bg-darksurface-100 dark:border-darksurface-400">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DropdownMenuSeparator className="dark:bg-darksurface-400" />
                <DropdownMenuItem className="dark:hover:bg-darksurface-400">
                    <Link href={`/logs/${encodeURIComponent(log.id)}`}>
                        View Details
                    </Link>
                </DropdownMenuItem>
                <DropdownMenuItem className="dark:hover:bg-darksurface-400">
                    View Exploit Code
                </DropdownMenuItem>
                <DropdownMenuSeparator className="dark:bg-darksurface-400" />
                <DropdownMenuItem className="dark:hover:bg-darksurface-400">Copy Request</DropdownMenuItem>
                <DropdownMenuItem className="dark:hover:bg-darksurface-400">Copy Response</DropdownMenuItem>
                <DropdownMenuItem className="dark:hover:bg-darksurface-400">Copy Exploit Code</DropdownMenuItem>
                <DropdownMenuSeparator className="dark:bg-darksurface-400" />
                <DropdownMenuItem className="bg-red-600 dark:hover:bg-red-600/[0.8]">
                    Delete Log
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          )
        },
      },
]