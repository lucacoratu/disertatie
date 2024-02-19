import { useTheme } from "next-themes";
import { FC } from "react";

import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
  } from "@/components/ui/tooltip"
import { rule } from "postcss";

type RuleFindingProps = {
    finding: RuleFindingData | null | undefined,
    matchedString: string | undefined,
}

const RuleFinding: FC<RuleFindingProps> = ({finding, matchedString}): JSX.Element => {
    const severityColors = ['bg-[#fdc500]', 'bg-[#fd8c00]', 'bg-[#dc0000]', 'bg-[#780000]'];
    const severityTextColors = ['text-[#fdc500]', 'text-[#fd8c00]', 'text-[#dc0000]', 'text-[#780000]'];
    const severityNames = ['Low', 'Medium', 'High', 'Critical'];

    return (
        <div key={finding?.id} className="overflow-hidden">
            {
                finding && finding.id != '' && 
                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger>
                            <div className="flex flex-row gap-2 items-center">
                                <div className={"w-5 h-5 rounded-full " + severityColors[finding.severity]}/>
                                <div>{finding.ruleId},</div>
                                <div>Severity: {severityNames[finding.severity]},</div>
                                <div>Classification: {finding.classification.toUpperCase()}</div>
                            </div>
                        </TooltipTrigger>
                        <TooltipContent>
                            <div className="flex flex-col gap-1 items-center">
                                <p>{finding.ruleDescription}</p>
                                <p>Classified as: {finding.classification.toUpperCase()}</p>
                                <p>Severity: <span className={severityTextColors[finding.severity]}>{severityNames[finding.severity]}</span> - Detected by: {finding.ruleName} (Id: {finding.ruleId})</p>
                                <p>Matched on line {finding.line + 1}, position {finding.lineIndex} {matchedString != "" ? "- String: " + matchedString : ""}</p>
                                <p>{finding.matchedBodyHash != "" ? "Matched on body hash (Algorithm: "+ finding.matchedBodyHashAlg + "): " + finding.matchedBodyHash : ""}</p>
                            </div>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            }
        </div>
    );
}

export default RuleFinding;