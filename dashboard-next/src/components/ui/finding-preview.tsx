import { useTheme } from "next-themes";
import { FC } from "react";

import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
  } from "@/components/ui/tooltip"

type FindingPreviewProps = {
    finding: Finding,
    findingsClassificationString: FindingClassificationString[],
}

const FindingPreview: FC<FindingPreviewProps> = ({finding, findingsClassificationString}): JSX.Element => {
    const severityColors = ['bg-[#fdc500]', 'bg-[#fd8c00]', 'bg-[#dc0000]', 'bg-[#780000]'];
    const severityTextColors = ['text-[#fdc500]', 'text-[#fd8c00]', 'text-[#dc0000]', 'text-[#780000]'];

    // if(finding.request.id != "") {
    //     console.log(finding);
    // }
    const severityNames = ['Low', 'Medium', 'High', 'Critical'];

    const requestClassificationString = findingsClassificationString.filter((classificationString) => classificationString.intFormat == finding.request.classification).at(0)?.stringFormat;
    const requestClassificationDescription = findingsClassificationString.filter((classificationString) => classificationString.intFormat == finding.request.classification).at(0)?.description;
    const responseClassificationString = findingsClassificationString.filter((classificationString) => classificationString.intFormat == finding.response.classification).at(0)?.stringFormat;
    const responseClassificationDescription = findingsClassificationString.filter((classificationString) => classificationString.intFormat == finding.response.classification).at(0)?.description;

    return (
        <div className="overflow-hidden">
            {
                finding.request.id != '' && 
                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger>
                            <div className={"p-1 rounded w-fit text-center " + severityColors[finding.request.severity]}>
                                {requestClassificationString}
                            </div>
                        </TooltipTrigger>
                        <TooltipContent>
                            <div className="flex flex-col gap-1 items-center">
                                <p>{requestClassificationDescription}</p>
                                <p>Severity: <span className={severityTextColors[finding.request.severity]}>{severityNames[finding.request.severity]}</span> - Detected by: {finding.request.validatorName}</p>
                            </div>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            }

            {
                finding.response.id != '' && 
                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger>
                            <div className={"p-1 rounded w-fit text-center " + severityColors[finding.response.severity]}>
                                {responseClassificationString}
                            </div>
                        </TooltipTrigger>
                        <TooltipContent>
                            <div className="flex flex-col gap-1 items-center">
                                <p>{responseClassificationDescription}</p>
                                <p>Severity: <span className={severityColors[finding.response.severity]}>{severityNames[finding.response.severity]}</span> - Detected by: {finding.response.validatorName}</p>
                            </div>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            }
        </div>
    );
}

export default FindingPreview;