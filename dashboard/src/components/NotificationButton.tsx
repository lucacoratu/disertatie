import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"

import { Bell } from "lucide-react";
import WebSocketConnection from "@/types/websocket";
import { useEffect, useState } from "react";
import {WSNotification, WsRuleDetectionAlert} from "@/types/websocket_types";
import { ScrollArea } from "@/components/ui/scroll-area"
import {DisplayAgentConnectedToast, DisplayAgentDisconnectedToast} from "@/lib/utils";

const NotificationButton = () => {
    const [notifications, setNotifications] = useState<JSX.Element[]>([]);
    const [numberNotifications, setNumberNotifications] = useState(0);

    useEffect(() => {
        //Connect to the websocket of the API
        const wsConnection: WebSocketConnection = WebSocketConnection.getInstance();

        wsConnection.handleNotification = (message: string) => {
            const notif: WSNotification = JSON.parse(message);
            console.log(notif);
            const updateNotifications = [
                // copy the current notifications state
                ...notifications,
                // now you can add a new object to add to the array
                <div key={notifications.length + 1}>{notif.message}</div>
            ];
            setNotifications(updateNotifications);
            setNumberNotifications(updateNotifications.length);
        };

        wsConnection.handleAgentDisconnectNotificationCallback = (message: string) => {
            const notif: WSNotification = JSON.parse(message);
            DisplayAgentDisconnectedToast(notif.agentId);
        };

        wsConnection.handleAgentConnectNotificationCallback = (message: string) => {
            const notif: WSNotification = JSON.parse(message);
            DisplayAgentConnectedToast(notif.agentId);
        };

        wsConnection.handleAgentRuleDetectionAlert = (message: string) => {
            const alert: WsRuleDetectionAlert = JSON.parse(message);

            const detectionTime: string = new Date(alert.timestamp * 1000).toLocaleString();

            const updateNotifications = [
                // copy the current notifications state
                ...notifications,
                // now you can add a new object to add to the array
                <div key={notifications.length + 1} className="flex flex-col bg-red-900 rounded b-1 p-2">
                    <div className="flex flex-col items-start">
                        <p className="text-base mb-2">Rule detection alert - {alert.severity.toUpperCase()}</p>
                        <div className="text-sm">Agent: {alert.agentId}</div>
                        <div className="text-sm">Rule: {alert.ruleName}</div>
                        <div className="text-sm">Classification: {alert.classification}</div>
                    </div>
                    <div className="mt-2 flex flex-row justify-between">
                        <div className="flex flex-row gap-1">
                            <div className="text-xs">Ignore</div>
                            <div className="text-xs">View</div>
                        </div>
                        <div className="text-xs">{detectionTime}</div>
                    </div>
                </div>
            ];
            setNotifications(updateNotifications);
            setNumberNotifications(updateNotifications.length);
        }
    }, [])
    

    return (
        <Popover>
            <PopoverTrigger className="h-6 w-6 flex flex-row items-center">
                <div className="flex flex-row items-center">
                    <Bell className="h-6 w-6"/>
                    {numberNotifications !== 0 && <div className="h-4 w-4 self-start -mt-1 -ml-2 text-xs text-center bg-red-500 rounded-full">{numberNotifications}</div>}
                </div>
            </PopoverTrigger>
            <PopoverContent className="w-[400px] max-h-[600px]">
                <ScrollArea className="w-full max-h-[550px] p-2 flex flex-col gap-1">
                    <div className="flex flex-col gap-2">
                        {notifications.map((notification, index) => {
                            return <div key={index}>{notification}</div>
                        })}
                    </div>
                </ScrollArea>
            </PopoverContent>
        </Popover>
    );
}

export default NotificationButton;