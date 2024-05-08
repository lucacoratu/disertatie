import { FC } from "react";

import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
  } from "@/components/ui/tooltip"

type FindingPreviewProps = {
    ruleFinding: RuleFinding,
}

const RuleFindingPreview: FC<FindingPreviewProps> = ({ruleFinding}): JSX.Element => {
    const severityColors = ['bg-[#fdc500]', 'bg-[#fd8c00]', 'bg-[#dc0000]', 'bg-[#780000]'];
    const severityTextColors = ['text-[#fdc500]', 'text-[#fd8c00]', 'text-[#dc0000]', 'text-[#780000]'];

    const severityNames = ['Low', 'Medium', 'High', 'Critical'];

    return (
        <div className="overflow-hidden" key={ruleFinding.request?.id || ruleFinding.response?.id}>
            {
                ruleFinding.request && ruleFinding.request?.id != '' && 
                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger>
                            <div className={"w-fit min-w-fit truncate ... overflow-hidden p-1 rounded font-bold text-center " + severityColors[ruleFinding.request?.severity || 0]}>
                                {ruleFinding.request?.classification.toUpperCase()}
                            </div>
                        </TooltipTrigger>
                        <TooltipContent>
                            <div className="flex flex-col gap-1 items-center">
                                <p>{ruleFinding.request?.ruleDescription}</p>
                                <p>Severity: <span className={severityTextColors[ruleFinding.request?.severity || 0]}>{severityNames[ruleFinding.request?.severity || 0]}</span> - Detected by: {ruleFinding.request?.ruleId}</p>
                            </div>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            }

            {
                ruleFinding.response && ruleFinding.response?.id != '' && 
                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger>
                            <div className={"truncate ... overflow-hidden p-1 rounded w-fit min-w-fit font-bold text-center " + severityColors[ruleFinding.response?.severity || 0]}>
                                {ruleFinding.response?.classification.toUpperCase()}
                            </div>
                        </TooltipTrigger>
                        <TooltipContent>
                            <div className="flex flex-col gap-1 items-center">
                                <p>{ruleFinding.response?.ruleDescription}</p>
                                <p>Severity: <span className={severityTextColors[ruleFinding.response?.severity || 0]}>{severityNames[ruleFinding.response?.severity || 0]}</span> - Detected by: {ruleFinding.response?.ruleId}</p>
                            </div>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            }
        </div>
    );
}

export default RuleFindingPreview;