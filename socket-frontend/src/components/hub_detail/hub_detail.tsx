import { Table, Tag } from "antd";
import { HubDetail } from "../../types/connection_detail";

import "./styles.css"
import "../../index.css"
import Column from "antd/es/table/Column";

const columns = [
    {
        title: 'CID',
        dataIndex: 'connection_id',
        key: 'connection_id',
    },
    {
        title: 'Status',
        dataIndex: 'status',
        key: 'status',
    },
    {
        title: 'Connected at',
        dataIndex: 'connected_at',
        key: 'connected_at',
    },
    {
        title: 'Disconnected at',
        dataIndex: 'disconnected_at',
        key: 'disconnected_at',
    }
];

function getStatus(status: number): string {
    return status === 0 ? "ALIVE" : "DEAD"
}

export default function HubDetailComponent(hubDetails: HubDetail) {
    return <div className="active-connection-box">
        <h3 className="active-connections">Active connections {hubDetails.active_connections}</h3>
        <div className="hub-detail">
            <Table
                rowKey={conn => conn.connection_id}
                dataSource={hubDetails.connections}>
                {columns.map((s) => <Column title={s.title} dataIndex={s.dataIndex} key={s.key}></Column>)}
                <Column
                    title="Topics"
                    dataIndex="subscribed_topics"
                    key="subscribed_topics"
                    render={(tags: string[]) => (
                        <>
                            {tags.map((tag) => (
                                <Tag color="blue" key={tag}>
                                    {tag}
                                </Tag>
                            ))}
                        </>
                    )}
                />
            </Table>
        </div>
    </div>
}