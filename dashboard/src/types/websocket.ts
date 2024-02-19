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

    public handleAgentDisconnectNotificationCallbacks: handleNotificationCallback[];
    private handleAgentConnectNotificationCallbacks: handleNotificationCallback[];

    private constructor() { 
        //Initialize the websocket connection
        this.conn = new WebSocket(constants.apiWebSocketURL);
        this.conn.onopen = () => { this.connOpen = true; }
        this.conn.onclose = () => { this.connOpen = false; }
        this.conn.onmessage = (event: MessageEvent) => this.handleMessage(event);

        this.handleNotification = (text: string) => {};
        this.handleAgentStatusResponse = (text: string) => {};
        this.handleAgentRuleDetectionAlert = (text: string) => {};

        this.handleAgentDisconnectNotificationCallbacks = [];
        this.handleAgentConnectNotificationCallbacks = [];
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
                    for(var i = 0; i < this.handleAgentDisconnectNotificationCallbacks.length; i++) {
                        this.handleAgentDisconnectNotificationCallbacks[i](content);
                    }
                    break;
                case WsMessageTypes.WsAgentConnectedNotification:
                    for(var i = 0; i < this.handleAgentConnectNotificationCallbacks.length; i++) {
                        this.handleAgentConnectNotificationCallbacks[i](content);
                    }
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

    public addAgentConnectedCallback(callback: handleNotificationCallback) {
        //Check if the function does not already exist in the list of callbacks
        if (this.handleAgentConnectNotificationCallbacks.indexOf(callback) != -1)
            return;

        this.handleAgentConnectNotificationCallbacks.push(callback);
    }

    // public deleteAgentConnectedCallback(callback: handleNotificationCallback) {
    //     for(let i = 0; i < this.handleAgentConnectNotificationCallbacks.length; i++) {
    //         if (this.handleAgentConnectNotificationCallbacks[i] === callback) {
    //             this.handleAgentConnectNotificationCallbacks = [...this.handleAgentConnectNotificationCallbacks.slice(0, i - 1), ...this.handleAgentConnectNotificationCallbacks.slice(i + 1, Infinity)];
    //             return;
    //         }
    //     }
    // }

    public addAgentDisconnectedCallback(callback: handleNotificationCallback) {
        //Check if the function does not already exist in the list of callbacks
        if (this.handleAgentDisconnectNotificationCallbacks.indexOf(callback) != -1)
            return;

        this.handleAgentDisconnectNotificationCallbacks.push(callback);
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