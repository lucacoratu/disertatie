import { FC } from "react";

const OsCard: FC<OSProps> = ({os}): JSX.Element => {
    //Check if the operating system is Windows and return the coresponding UI element
    if(os?.toLowerCase() === "windows") {
        return (
            <div className="min-w-9 p-1 b-1 rounded-lg flex items-center bg-blue-800">
                <img className='machine-os-image' src='/icons8-windows-30.png'>
                </img>
            </div>
        );
    }

    if(os?.toLowerCase() === "linux") {
        //If the os is linux
        return (
            <div className="min-w-9 p-1 b-1 rounded-lg flex items-center bg-orange-600">
                <img className='machine-os-image' src='/icons8-linux-30.png'>
                </img>
            </div>
        );
    }

    return (
        <div className="min-w-9 p-1 b-1 rounded-lg flex items-center bg-orange-600">
            <img className='machine-os-image' src='/icons8-linux-30.png'>
            </img>
        </div>
    )
}

export default OsCard;