import { FC } from "react";
import { constants } from "@/app/constants";
import Finding from "@/components/ui/finding";
import { request } from "http";

type FindingCardProps = {
    findings: Finding[];
    findingsClassificationString: FindingClassificationString[],
}

const FindingCard: FC<FindingCardProps> = ({findings, findingsClassificationString}): JSX.Element => {
    //Separate request findings from response findings
    const requestFindings: (FindingData | null)[] = findings?.map((finding) => {
        return finding.request.id !== "" ? finding.request : null;
    }).filter((finding) => finding !== null);

    //console.log(requestFindings);

    const responseFindings: (FindingData | null)[] = findings?.map((finding) => {
        return finding.response.id !== "" ? finding.response : null;
    }).filter((finding) => finding !== null);

    //console.log(responseFindings);

    return (
        <div className="flex flex-wrap flex-row gap-3 grow min-h-[200px] justify-center">
            {/* Request findings div*/}
            <div className="flex flex-col gap-1 min-w-[450px] grow w-1/3">
                <div>
                    <h2 className="text-lg text-gray-400">Request</h2>
                </div>
                <div>
                    {requestFindings && requestFindings.length !== 0 && requestFindings.map((finding) => {
                        return (
                            <Finding key={finding?.id} finding={finding} findingsClassificationString={findingsClassificationString} matchedString={finding?.matchedString}/>
                        );
                    })}

                    {requestFindings && requestFindings.length === 0 && 
                        <div className="flex items-center text-center text-gray-600">
                            No request findings
                        </div>
                    }
                </div>
            </div>
            {/* Response findings div */}
            <div className="flex flex-col gap-1 min-w-[450px] grow w-1/3">
                <div>
                    <h2 className="text-lg text-gray-400">Response</h2>
                </div>
                <div>
                    {responseFindings && responseFindings.length !== 0 && responseFindings.map((finding) => {
                        return (
                            <Finding key={finding?.id} finding={finding} findingsClassificationString={findingsClassificationString} matchedString={finding?.matchedString}/>
                        );
                    })}

                    {responseFindings && responseFindings.length == 0 && 
                        <div className="flex items-center text-center text-gray-600">
                            No response findings
                        </div>
                    }
                </div>
            </div>
        </div>
    );
}

export default FindingCard;