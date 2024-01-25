import { Button, Table } from "antd";
import { ConnectionDetail } from "../../types/connection_detail";
import "./styles.css";

function getStatus(status: number): string {
    return status === 0 ? "ALIVE" : "DEAD"
}

export default function ConnectionDetailComponent(connection: ConnectionDetail) {
    return <div className="connection-detail">
        <h3>Connection id {connection.connection_id}</h3>
        <h3>Connection status {getStatus(connection.status)}</h3>
        <h3>Connected at {connection.connected_at}</h3>
        <h3>Disconnected at {connection.disconnected_at}</h3>
        Suscribed topics {connection.subscribed_topics.map((topic) => <h4>{topic}</h4>)}
    </div>
}