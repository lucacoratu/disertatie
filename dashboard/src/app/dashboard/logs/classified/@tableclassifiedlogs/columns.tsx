"use client";

import RuleFindingPreview from "@/components/ui/rulefinding-preview";
import FindingPreview from "@/components/ui/finding-preview";
import { Checkbox } from "@/components/ui/checkbox";

import { ColumnDef } from "@tanstack/react-table";
import { Power, Square, MoreHorizontal, Repeat, ArrowUpDown, Pencil, Trash } from "lucide-react";
 
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { DataTableColumnHeader } from "@/components/table/column-header";
import Link from "next/link";

//Create another type which will define the columns of the table
export type LogColumn = {
    id: string,
    agentId: string,
    remoteip: string,
    requestmethod: string,
    url: string,
    response: string,
    timestamp: string,
    findings: Finding[],
    rulefindings: RuleFinding[],
    findingsClassificationString: FindingClassificationString[],
    nextPage: string,
}

export const columns: ColumnDef<LogColumn>[] = [
  {
    id: "select",
    header: ({ table }) => (
    <Checkbox
      checked={
        table.getIsAllPageRowsSelected() ||
        (table.getIsSomePageRowsSelected() && "indeterminate")
      }
      onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
      aria-label="Select all"
      className="mx-auto"
    />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
        className="mx-auto"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "remoteip",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Remote IP" />
    ),
    cell: ({ row }) => {
      return (
        <p className="text-sm font-bold">{row.getValue('remoteip')}</p>
      )
    }
  },
  {
    accessorKey: "requestmethod",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Request Method" />
    ),
  },
  {
    accessorKey: "url",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="URL" />
    ),
    cell: ({ row }) => {
      return (
        <div className="truncate text-nowrap w-full max-w-[400px]">
          {row.getValue('url')}
        </div>
      )
    }
  },
  {
    accessorKey: "response",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Response" />
    ),
  },
  {
    accessorKey: "timestamp",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="timestamp" />
    ),
  },
  {
    accessorKey: "findings",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Findings" />
    ),
    cell: ({ row }) => {
      const log = row.original;
      //console.log(log);
      if(log.findings != null) {
        return (
            <div className="overflow-hidden flex flex-row gap-2 max-w-36 max-h-8">
              {
                log.findings.map((finding: Finding) => (
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
    accessorKey: "rulefindings",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Rule Findings" />
    ),
    cell: ({ row }) => {
      const log = row.original;
        if(log.rulefindings != null) {
          return (
              <div className="overflow-hidden flex flex-row gap-2 max-w-52 max-h-8">
                {log.rulefindings.map((ruleFinding: RuleFinding) => (
                    <RuleFindingPreview key={ruleFinding.request?.id} ruleFinding={ruleFinding}/>
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
            <Button variant="ghost" className="h-8 w-8 p-0 max-w-9">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="dark:bg-darksurface-100 dark:border-darksurface-400">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuSeparator className="dark:bg-darksurface-400" />
            <DropdownMenuItem className="dark:hover:bg-darksurface-400">
                <Link href={`/dashboard/logs/${encodeURIComponent(log.id)}`}>
                    View Details
                </Link>
            </DropdownMenuItem>
            <DropdownMenuItem className="dark:hover:bg-darksurface-400">
              <Link href={`/dashboard/logs/${encodeURIComponent(log.id)}/exploit`}>
                View Exploit Code
              </Link>
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
];