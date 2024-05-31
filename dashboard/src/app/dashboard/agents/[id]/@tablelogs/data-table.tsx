"use client"
 
import * as React from "react"
 
import { cn } from "@/lib/utils"
import { Check, ChevronDown } from "lucide-react";


import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
 


import {
  ColumnDef,
  ColumnFiltersState,
  VisibilityState,
  SortingState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table"
 
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"

import { Input } from "@/components/ui/input"
import { constants } from "@/app/constants";
import { DropdownMenuGroup, DropdownMenuItem } from "@radix-ui/react-dropdown-menu";
import { useCookies } from 'next-client-cookies';
 
interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  agentId: string
}

const tableColumns = [
  {
    value: "remoteip",
    label: "Remote IP",
  },
  {
    value: "requestmethod",
    label: "Request Method",
  },
  {
    value: "url",
    label: "URL",
    columnName: "url",
  },
  {
    value: "response",
    label: "Response",
  },
  {
    value: "timestamp",
    label: "Timestamp",
  },
]
 
export function DataTable<TData, TValue>({
  columns,
  data,
  agentId,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  )
  const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({})

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    getSortedRowModel: getSortedRowModel(),
    onColumnFiltersChange: setColumnFilters,
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
    },
  })
 
  const [open, setOpen] = React.useState(false);
  const [searchColumn, setSearchColumn] = React.useState("");

  const cookie = useCookies().get('session');

  return (
    <div>
      <div className="flex items-center py-4 gap-1">
          <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              role="combobox"
              aria-expanded={open}
              className="w-[200px] justify-between"
            >
              {searchColumn
                ? tableColumns.find((tableColumn) => tableColumn.value === searchColumn)?.label
                : "Select column..."}
              <ChevronDown/>
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-[200px] p-0">
            <Command>
              <CommandInput placeholder="Search columns..." className="h-9" />
              <CommandEmpty>No column found.</CommandEmpty>
              <CommandGroup>
                {tableColumns.map((tableColumn) => (
                  <CommandItem
                    key={tableColumn.value}
                    value={tableColumn.value}
                    onSelect={(currentValue) => {
                      setSearchColumn(currentValue === searchColumn ? "" : currentValue);
                      setOpen(false);
                    }}
                  >
                    {tableColumn.label}
                    <Check
                      className={cn(
                        "ml-auto h-4 w-4",
                        searchColumn === tableColumn.value ? "opacity-100" : "opacity-0"
                      )}
                    />
                  </CommandItem>
                ))}
              </CommandGroup>
            </Command>
          </PopoverContent>
        </Popover>

        <Input
          placeholder={`Filter ${searchColumn}...`}
          value={(table.getColumn(searchColumn)?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn(searchColumn)?.setFilterValue(event.target.value)
          }
          className="max-w-sm"
        />

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="ml-auto h-8" >
              Columns
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {table
              .getAllColumns()
              .filter(
                (column) => column.getCanHide()
              )
              .map((column) => {
                return (
                  <DropdownMenuCheckboxItem
                    key={column.id}
                    className="capitalize"
                    checked={column.getIsVisible()}
                    onCheckedChange={(value) =>
                      column.toggleVisibility(!!value)
                    }
                  >
                    {column.id}
                  </DropdownMenuCheckboxItem>
                )
              })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className="rounded-md border">
        <Table>
          <TableHeader className="bg-card">
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  )
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row, index) => {
                let styling = ""
                if(index % 2 == 1) {
                  styling = "dark:bg-darksurface-100";
                } else {
                  styling = "bg-card";
                }
                styling += " bg-card";
                return <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                  className={styling}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              })
            ) : (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center">
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
        {/* <div className="flex items-center justify-end space-x-2 py-4">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Next
          </Button>
        </div> */}
          <Pagination>
            <PaginationContent>
                <PaginationPrevious onClick={() => {
                  table.previousPage();
                }}/>
                {/* {[...Array(table.getPageCount())].map((page) => {
                    return <PaginationLink href="#">{page}</PaginationLink>;
                })} */}
                <PaginationEllipsis />
                <PaginationNext onClick={() => {
                  table.nextPage();
                }}/>
            </PaginationContent>
          </Pagination>
      </div>
  )
}