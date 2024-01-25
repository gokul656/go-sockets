export interface HubDetail {
    active_connections: number,
    connections: ConnectionDetail[]
}

export interface ConnectionDetail {
    connection_id: string,
    status: number,
    subscribed_topics: string[],
    connected_at: string,
    disconnected_at: string,
}

export enum ConnectionStatus {
    ALIVE, DEAD
}

export interface TopicData {
    topic: string,
    status: boolean
}