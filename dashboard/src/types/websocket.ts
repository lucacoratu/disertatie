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

    private constructor() { 
        //Initialize the websocket connection
        this.conn = new WebSocket(constants.apiWebSocketURL);
        this.conn.onopen = () => {this.connOpen = true;}
        this.conn.onmessage = (event: MessageEvent) => this.handleMessage(event);
        this.handleNotification = (text: string) => {};
        this.handleAgentStatusResponse = (text: string) => {};
    }

    private handleMessage(event: MessageEvent) {
        const message: WSMessage = JSON.parse(event.data);
        //console.log(message);
        //console.log(message.content);
        if (message) {
            switch(message.type) {
                case WsMessageTypes.WSNotification:
                    const content: string = JSON.stringify(message.data);
                    this.handleNotification(content);
                case WsMessageTypes.WsAgentStatusResponse:
                    const msgdata: string = JSON.stringify(message.data);
                    this.handleAgentStatusResponse(msgdata)
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