import { FC } from "react";
import { constants } from "@/app/constants";

type HTTPHighlighterProps = {
    log: LogFull
}

//Returns the line in a highlighted fashion based on the findings
function HighlightRequestLine(line: string, lineNumber: number, findings: Finding[]): string[] {
    //Initialize the new line
    let highlightedLine: string = '<p className="text-base line-clamp-1">';
    const highlights: string[] = [];
    //Loop through each finding
    findings?.map((finding) => {
        //Check if the finding is for this line
        if(finding.request.line == lineNumber && finding.request.id != "") {
            //The finding is for this line, so create the highlights
            // console.log(line.length);
            // console.log(finding.request.line);
            // console.log(finding.request.lineIndex, finding.request.lineIndex + finding.request.length);
            // console.log(line.slice(finding.request.lineIndex, finding.request.lineIndex + finding.request.length));
            let highlight: string = "<span className=" + constants.severityTextColors[finding.request.severity] + ">"+ line.slice(finding.request.lineIndex, finding.request.lineIndex + finding.request.length) + "</span>";
            console.log(highlight);
            //highlights.push(highlight);
        }
    });


    return highlights;
    // return (
    //     <>
    //     </>
    // );
}

const HTTPHighlighter: FC<HTTPHighlighterProps> = ({log}): JSX.Element => {
    //Get the request data from base64
    const logRequestRaw: string = atob(log?.request);
    //Get the response data from base64
    const logResponseRaw: string = atob(log?.response);

    //console.log(log);

    return (
        <div className="flex flex-wrap flex-row gap-5 justify-center">
            <div className="flex flex-col gap-0 min-w-[450px] grow w-1/3 p-4 rounded dark:bg-darksurface-100">
                <div className="text-center">
                    <h2>Request</h2>
                </div>
                <div className="flex flex-col gap-0">
                {logRequestRaw.split('\n').map((line, index) => {
                    const auxLine = line.trimEnd();
                    //console.log(log?.findings);
                    const out = HighlightRequestLine(line, index, log.findings);
                    //console.log(out);
                    return <p className="text-base line-clamp-1" key={index}>{auxLine}</p>;
                })}
                </div>
            </div>
            <div className="flex flex-col gap-0 min-w-[450px] grow w-1/3 p-4 rounded dark:bg-darksurface-100">
                <div className="text-center">
                    <h2>Response</h2>
                </div>
                <div>
                    {logResponseRaw.split('\n').map((line, index) => {
                        const auxLine = line.trimEnd();
                        //HighlightRequestLine(line, index, log.findings);
                        return <p className="text-base line-clamp-1" key={index}>{auxLine}</p>;
                    })}
                </div>
            </div>
        </div>
    );
}

export default HTTPHighlighter;