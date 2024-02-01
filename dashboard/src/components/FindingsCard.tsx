import { FC } from "react";
import { constants } from "@/app/constants";
import Finding from "@/components/ui/finding";

type FindingCardProps = {
    findings: Finding[];
    findingsClassificationString: FindingClassificationString[],
}

const FindingCard: FC<FindingCardProps> = ({findings, findingsClassificationString}): JSX.Element => {
    //Separate request findings from response findings
    const requestFindings: (FindingData | null)[] = findings.map((finding) => {
        return finding.request.id !== "" ? finding.request : null;
    }).filter((finding) => finding !== null);

    //console.log(requestFindings);

    const responseFindings: (FindingData | null)[] = findings.map((finding) => {
        return finding.response.id !== "" ? finding.response : null;
    }).filter((finding) => finding !== null);

    //console.log(responseFindings);

    return (
        <div className="flex flex-wrap flex-row gap-3 grow min-h-[200px] justify-center">
            {/* Request findings div*/}
            <div className="flex flex-col gap-1 min-w-[450px] grow w-1/3 p-4 rounded dark:bg-darksurface-100 dark:border-darksurface-400 border-2 dark:bg-darksurface-100 b-2 rounded-xl">
                <div>
                    <h2 className="text-xl">Request findings</h2>
                </div>
                <div>
                    {requestFindings.map((finding) => {
                        return (
                            <Finding key={finding?.id} finding={finding} findingsClassificationString={findingsClassificationString} matchedString={finding?.matchedString}/>
                        );
                    })}
                </div>
            </div>
            {/* Response findings div */}
            <div className="flex flex-col gap-1 min-w-[450px] grow w-1/3 p-4 rounded dark:bg-darksurface-100 dark:border-darksurface-400 border-2 dark:bg-darksurface-100 b-2 rounded-xl">
                <div>
                    <h2 className="text-xl">Response findings</h2>
                </div>
                <div>
                    {responseFindings.map((finding) => {
                        return (
                            <Finding key={finding?.id} finding={finding} findingsClassificationString={findingsClassificationString} matchedString={finding?.matchedString}/>
                        );
                    })}
                </div>
            </div>
        </div>
    );
}

export default FindingCard;