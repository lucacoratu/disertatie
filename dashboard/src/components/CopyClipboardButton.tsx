"use client";

import { FC } from "react";
import { toast } from "sonner";

import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
  } from "@/components/ui/tooltip";

import { Button } from "@/components/ui/button";

type CopyClipboardButtonProps = {
    text: string,
    tooltipText: string,
    toastText: string,
    children: JSX.Element,
}

function copyText(text: string, tooltipText: string){
    navigator.clipboard.writeText(text);
    toast(tooltipText, {
        description: "Copied text to clipboard",
        action: {
            label: "Ok",
            onClick: () => {},
        }
    });
}

const CopyClipboardButton: FC<CopyClipboardButtonProps> = ({text, tooltipText,toastText, children}): JSX.Element => {
    return (
        <TooltipProvider>
            <Tooltip>
                <TooltipTrigger onClick={() => copyText(text, toastText)}>
                    {children}
                </TooltipTrigger>
                <TooltipContent>
                    <p>{tooltipText}</p>
                </TooltipContent>
            </Tooltip>
        </TooltipProvider>
    );
}

export default CopyClipboardButton;