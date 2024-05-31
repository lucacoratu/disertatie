"use client";

import { FC } from "react";
import Link from "next/link";
import { useRouter } from 'next/navigation'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism';
import CopyCliboardButton from "@/components/CopyClipboardButton";
import { Copy } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area"
import {Button} from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

type clickHandler = { (agentId: string, format: string): void};

type ExportButtonProps = {
    agentId: string,
    clickHandler: clickHandler,
}

const ExportButton: FC<ExportButtonProps> = ({agentId, clickHandler}): JSX.Element => {
    return (
        <DropdownMenu>
        <DropdownMenuTrigger asChild>
            <Button className="ml-auto h-8">
                Export
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="center">
            <DropdownMenuLabel className="text-center h-8">Format</DropdownMenuLabel>
            <DropdownMenuSeparator className="bg-foreground/[0.3]"/>
            <DropdownMenuItem className="capitalize text-center hover:bg-accent h-8" onClick={async () => await clickHandler(agentId, "json")}>JSON</DropdownMenuItem>
            <DropdownMenuItem className="capitalize text-center hover:bg-accent h-8" onClick={async () => await clickHandler(agentId, "csv")}>CSV</DropdownMenuItem>
        </DropdownMenuContent>
    </DropdownMenu>
    );
}

export default ExportButton;