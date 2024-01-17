import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"

import { Bell } from "lucide-react";

const NotificationButton = () => {
    return (
        <Popover>
            <PopoverTrigger className="h-6 w-6 flex flex-row items-center">
                <Bell />
            </PopoverTrigger>
            <PopoverContent className="flex flex-col gap-1">
                <div>
                    <div>Notification 1</div>
                    <div>Notification 1</div>
                    <div>Notification 1</div>
                </div>
            </PopoverContent>
        </Popover>
    );
}

export default NotificationButton;