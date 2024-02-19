import { useTheme } from "next-themes";
import { FC } from "react";

import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
  } from "@/components/ui/tooltip"

type FindingProps = {
    finding: FindingData | null,
    findingsClassificationString: FindingClassificationString[],
    matchedString: string | undefined,
}

const Finding: FC<FindingProps> = ({finding, findingsClassificationString, matchedString}): JSX.Element => {
    const severityColors = ['bg-[#fdc500]', 'bg-[#fd8c00]', 'bg-[#dc0000]', 'bg-[#780000]'];
    const severityTextColors = ['text-[#fdc500]', 'text-[#fd8c00]', 'text-[#dc0000]', 'text-[#780000]'];
    const severityNames = ['Low', 'Medium', 'High', 'Critical'];

    const classificationString = findingsClassificationString.filter((classificationString) => classificationString.intFormat == finding?.classification).at(0)?.stringFormat;
    const classificationDescription = findingsClassificationString.filter((classificationString) => classificationString.intFormat == finding?.classification).at(0)?.description;

    return (
        <div key={finding?.id} className="overflow-hidden">
            {
                finding && finding.id != '' && 
                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger>
                            <div className="flex flex-row gap-2 items-center">
                                <div className={"w-5 h-5 rounded-full " + severityColors[finding.severity]}/>
                                <div>{classificationString}</div>
                            </div>
                        </TooltipTrigger>
                        <TooltipContent>
                            <div className="flex flex-col gap-1 items-center">
                                <p>{classificationDescription}</p>
                                <p>Severity: <span className={severityTextColors[finding.severity]}>{severityNames[finding.severity]}</span> - Detected by: {finding.validatorName}</p>
                                <p>Matched on line {finding.line + 1}, position {finding.lineIndex} - String: {matchedString}</p>
                            </div>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            }
        </div>
    );
}

export default Finding;