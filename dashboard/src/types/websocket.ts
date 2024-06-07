import { constants } from "@/app/constants";
import {WsMessageTypes, WSMessage, WSAgentStatusRequest} from "@/types/websocket_types";

type handleNotificationCallback = { (text: string): void};
type handleAgentStatusResponseCallback = { (text: string): void};

class WebSocketConnection {
    private static instance: WebSocketConnection;

    private conn: WebSocket;
    private connOpen: boolean = false; 
    public handleNotification: handleNotificationCallback;
    public handleAgentStatusResponse: handleAgentStatusResponseCallback;
    public handleAgentRuleDetectionAlert: handleNotificationCallback;

    public handleAgentDisconnectNotificationCallback: handleNotificationCallback;
    public handleAgentConnectNotificationCallback: handleNotificationCallback;
    public handleAgentConnectStatus: handleNotificationCallback;
    public handleAgentDisconnectStatus: handleNotificationCallback;

    private constructor() { 
        //Initialize the websocket connection
        this.conn = new WebSocket(constants.apiWebSocketURL);
        this.conn.onopen = () => { this.connOpen = true; }
        this.conn.onclose = () => { console.log("Websocket connection closed"); this.connOpen = false; }
        this.conn.onmessage = (event: MessageEvent) => this.handleMessage(event);

        this.handleNotification = (text: string) => {};
        this.handleAgentStatusResponse = (text: string) => {};
        this.handleAgentRuleDetectionAlert = (text: string) => {};

        this.handleAgentDisconnectNotificationCallback = (text: string) => {};
        this.handleAgentConnectNotificationCallback = (text: string) => {};
        this.handleAgentConnectStatus = (text: string) => {};
        this.handleAgentDisconnectStatus = (text: string) => {};
    }

    private handleMessage(event: MessageEvent) {
        const message: WSMessage = JSON.parse(event.data);
        console.log(message);
        if (message) {
            const content: string = JSON.stringify(message.data);
            switch(message.type) {
                case WsMessageTypes.WSNotification:
                    this.handleNotification(content);
                    break;
                case WsMessageTypes.WsAgentStatusResponse:
                    this.handleAgentStatusResponse(content);
                    break;
                case WsMessageTypes.WsAgentDisconnectedNotification:
                    this.handleAgentDisconnectNotificationCallback(content);
                    this.handleAgentDisconnectStatus(content);
                    break;
                case WsMessageTypes.WsAgentConnectedNotification:
                    this.handleAgentConnectNotificationCallback(content);
                    this.handleAgentConnectStatus(content);
                    break;
                case WsMessageTypes.WsAgentRuleDetectionAlert:
                    this.handleAgentRuleDetectionAlert(content);
                    break;
                default:
                    break
            }
        }
    }

    public static getInstance(): WebSocketConnection {
        if (!WebSocketConnection.instance) {
            WebSocketConnection.instance = new WebSocketConnection();
        }

        return WebSocketConnection.instance;
    }

    public getAgentStatus(agent_id: string): void {
        //Create the request message
        const request: WSAgentStatusRequest = {agentId: agent_id}
        const wsMessage: WSMessage = {type: WsMessageTypes.WsAgentStatusRequest, data: request}
        if (this.connOpen) {
            //Send a message on the websocket
            this.conn.send(JSON.stringify(wsMessage));
        }
    }
}

export default WebSocketConnection;