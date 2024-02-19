export const enum WsMessageTypes {
    WsError = -1,
    WSNotification = 1,
    WsAgentStatusRequest = 2,
    WsAgentStatusResponse = 3,
    WsAgentDisconnectedNotification = 4,
    WsAgentConnectedNotification = 5,
    WsAgentRuleDetectionAlert = 6,
}

export type WSMessage = {
    type: number,
    data: any,
}

export type WSNotification = {
    agentId: string,
    message: string
}

export type WSAgentStatusRequest = {
	agentId: string, 
}

export type WSAgentStatusResponse = {
	agentId: string, 
	status:  string
}

export type WsRuleDetectionAlert = {
    agentId:         string,
	ruleId:          string, 
	ruleName:        string,
	ruleDescription: string,
	classification:  string,
	severity:        string,
	timestamp:       number, 
}