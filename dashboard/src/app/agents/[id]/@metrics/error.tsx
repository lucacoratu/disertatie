"use client"

export default async function MetricsPageError({ params }: { params: { id: string } }) {
    const agentId: string = params.id; 
    
    return (
        <div className="text-center">
            Error
        </div>
    );
}