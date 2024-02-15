import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"

import { Bell } from "lucide-react";
import WebSocketConnection from "@/types/websocket";
import { useState } from "react";
import { toast } from "sonner";
import {WSNotification} from "@/types/websocket_types";

const NotificationButton = () => {
    const [notifications, setNotifications] = useState<string[]>([]);
    const [numberNotifications, setNumberNotifications] = useState(0);

    //Connect to the websocket of the API
    const wsConnection: WebSocketConnection = WebSocketConnection.getInstance();

    wsConnection.handleNotification = (message: string) => {
        const notif: WSNotification = JSON.parse(message);
        //console.log(notif);
        const updateNotifications = [
            // copy the current users state
            ...notifications,
            // now you can add a new object to add to the array
            notif.message
          ];
        setNotifications(updateNotifications);
        setNumberNotifications(updateNotifications.length);
    };

    return (
        <Popover>
            <PopoverTrigger className="h-6 w-6 flex flex-row items-center">
                <div className="flex flex-row items-center">
                    <Bell className="h-6 w-6"/>
                    {numberNotifications !== 0 && <div className="h-4 w-4 self-start -mt-1 -ml-2 text-xs text-center bg-red-500 rounded-full">{numberNotifications}</div>}
                </div>
            </PopoverTrigger>
            <PopoverContent className="flex flex-col gap-1">
                <div>
                    {notifications.map((notification, index) => {
                        return <div key={index}>{notification}</div>
                    })}
                </div>
            </PopoverContent>
        </Popover>
    );
}

export default NotificationButton;