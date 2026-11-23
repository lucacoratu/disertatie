type Constants = {
    appName: string,
    appVersion: string,
    apiBaseURL: string,
    apiWebSocketURL: string,
    severityColors: string[],
    severityTextColors: string[],
    severityNames: string[],
}

export const constants: Constants = {
    appName: "Dashboard",
    appVersion: "1.0.0",
    //Used by the server side of the dashboard (set at build time, defaults to a local API)
    apiBaseURL: process.env.NEXT_PUBLIC_API_URL || "http://127.0.0.1:8081/api/v1",
    //Used by the browser, so it must be an address of the API reachable by the users of the dashboard
    apiWebSocketURL: process.env.NEXT_PUBLIC_API_WS_URL || "ws://127.0.0.1:8081/api/v1/ws",
    severityColors: ['bg-[#fdc500]', 'bg-[#fd8c00]', 'bg-[#dc0000]', 'bg-[#780000]'],
    severityTextColors: ['text-[#fdc500]', 'text-[#fd8c00]', 'text-[#dc0000]', 'text-[#780000]'],
    severityNames: ['Low', 'Medium', 'High', 'Critical'],
} 

export const protectedRoutes = ["/dashboard"];
export const authRoutes = ["/", "/login"];
export const publicRoutes = [];