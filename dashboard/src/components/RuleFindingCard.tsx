import { FC } from "react";
import { constants } from "@/app/constants";
import RuleFinding from "@/components/ui/rulefinding";

type FindingCardProps = {
    findings: RuleFinding[];
}

const RuleFindingCard: FC<FindingCardProps> = ({findings}): JSX.Element => {
    //Separate request findings from response findings
    const requestFindings: (RuleFindingData | null | undefined)[] = findings?.map((finding) => {
        return finding.request?.id !== "" ? finding.request : null;
    }).filter((finding) => finding !== null);

    //console.log(requestFindings);

    const responseFindings: (RuleFindingData | null | undefined)[] = findings?.map((finding) => {
        return finding.response?.id !== "" ? finding.response : null;
    }).filter((finding) => finding !== null);

    //console.log(responseFindings);

    return (
        <div className="flex flex-wrap flex-row gap-3 grow min-h-[200px] justify-center">
            {/* Request findings div*/}
            <div className="flex flex-col gap-1 min-w-[450px] grow w-1/3 p-4 rounded dark:bg-darksurface-100 dark:border-darksurface-400 border-2 dark:bg-darksurface-100 b-2 rounded-xl">
                <div>
                    <h2 className="text-xl">Request rule findings</h2>
                </div>
                <div className="flex flex-row gap-1 flex-wrap">
                    {requestFindings?.map((finding) => {
                        return (
                            <RuleFinding key={finding?.id} finding={finding} matchedString={finding?.matchedString}/>
                        );
                    })}
                </div>
            </div>
            {/* Response findings div */}
            <div className="flex flex-col gap-1 min-w-[450px] grow w-1/3 p-4 rounded dark:bg-darksurface-100 dark:border-darksurface-400 border-2 dark:bg-darksurface-100 b-2 rounded-xl">
                <div>
                    <h2 className="text-xl">Response rule findings</h2>
                </div>
                <div className="flex flex-row gap-1 flex-wrap">
                    {responseFindings?.map((finding) => {
                        return (
                            <RuleFinding key={finding?.id} finding={finding} matchedString={finding?.matchedString}/>
                        );
                    })}
                    
                    {responseFindings.length == 0 && 
                        <div className="flex items-center text-center text-accent">
                            No response rule findings
                        </div>
                    }
                </div>
            </div>
        </div>
    );
}

export default RuleFindingCard;