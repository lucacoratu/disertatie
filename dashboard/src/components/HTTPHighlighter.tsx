"use client"

import { FC } from "react";
import { constants } from "@/app/constants";
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { Copy, Scroll } from "lucide-react";

import CopyCliboardButton from "@/components/CopyClipboardButton";

import { Button } from "@/components/ui/button";

import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
  } from "@/components/ui/tooltip"

import { useTheme } from "next-themes";

import { ScrollArea } from "@/components/ui/scroll-area"

type HTTPHighlighterProps = {
    log: LogFull
}

function SeparateHeadersAndBody(data: string): string[] {
    //Get the body of the response
    const dataLines: string[] = data.split("\n\n");
    var dataParts: string[][] = dataLines.map((line) => line.split("\n"));
    if(dataParts.length == 0) {
        return ["", ""];
    }
    var dataHeaders: string = ""; 
    var dataBody: string = "" ;
    if(dataParts.length >= 1) {
        dataHeaders = dataParts[0].join("\n");
    }
    if (dataParts.length == 2) {
        dataBody = dataParts[1].join("\n");
    }
    return [dataHeaders, dataBody];
}

const HTTPHighlighter: FC<HTTPHighlighterProps> = ({log}): JSX.Element => {
    //Get the request data from base64
    const logRequestRaw: string = atob(log?.request);
    //Get the response data from base64
    const logResponseRaw: string = atob(log?.response);

    //Separate the headers from the body of the request
    const requestValues = SeparateHeadersAndBody(logRequestRaw);
    const requestHeaders = requestValues[0]; 
    const requestBody = requestValues[1]; 

    //Separate the headers from the body of the response
    const responseValues = SeparateHeadersAndBody(logResponseRaw);
    const responseHeaders = responseValues[0]; 
    const responseBody = responseValues[1];

    const {theme, setTheme} = useTheme();

    return (
        <div className="flex flex-wrap flex-row gap-3 justify-center">
            <div className="flex flex-col gap-0 h-[500px] min-w-[450px] grow w-1/3 p-4 rounded dark:bg-darksurface-100 dark:border-darksurface-400 border-2 dark:bg-darksurface-100 b-2 rounded-xl">
                <div className="flex flex-row items-center justify-between gap-10">
                    <h2 className="text-xl">Request</h2>
                    <CopyCliboardButton text={logRequestRaw} toastText="Request copied to clipboard" tooltipText="Copy to clipboard">
                        <Copy className="w-4 h-4"/>
                    </CopyCliboardButton>
                </div>
                <ScrollArea>
                    {requestHeaders != "" &&
                        <SyntaxHighlighter language="http" style={theme === "light" ? oneLight : oneDark} showLineNumbers>
                            {requestHeaders}
                        </SyntaxHighlighter>
                    }
                    {requestBody != "" && 
                        <SyntaxHighlighter language="json" style={theme === "light" ? oneLight : oneDark} showLineNumbers>
                            {requestBody}
                        </SyntaxHighlighter>
                    }
                </ScrollArea>
            </div>
            <div className="flex flex-col gap-0 h-[500px] min-w-[450px] grow w-1/3 p-4 rounded dark:bg-darksurface-100 dark:border-darksurface-400 border-2 dark:bg-darksurface-100 b-2 rounded-xl">
                <div className="flex flex-row items-center justify-between gap-10">
                    <h2 className="text-xl">Response</h2>
                    <CopyCliboardButton text={logResponseRaw} toastText="Response copied to clipboard" tooltipText="Copy to clipboard" >
                        <Copy className="w-4 h-4"/>
                    </CopyCliboardButton>
                </div>
                <ScrollArea>
                    {responseHeaders != "" &&
                        <SyntaxHighlighter language="http" style={theme === "light" ? oneLight : oneDark} showLineNumbers>
                            {responseHeaders}
                        </SyntaxHighlighter>
                    }
                    {responseBody != "" && 
                        <SyntaxHighlighter language="html" style={theme === "light" ? oneLight : oneDark} showLineNumbers>
                            {responseBody}
                        </SyntaxHighlighter>
                    }
                </ScrollArea>
            </div>
        </div>
    );
}

export default HTTPHighlighter;