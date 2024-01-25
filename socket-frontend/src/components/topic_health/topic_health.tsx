import { useEffect, useState } from "react"
import "./styles.css"
import { getAvailableTopics } from "../../service/api_service"
import { HealthStatus } from "../status/status"
import { TopicData } from "../../types/connection_detail"

export default function TopicHealth() {

    const [isLoading, setLoading] = useState(true)
    const [availableTopics, setAvailableTopics] = useState<TopicData[]>([])

    useEffect(() => {
        const availableTopic = (async () => {
            const res = await getAvailableTopics();
            setAvailableTopics(res)
            setLoading(false)
        })
        availableTopic()
    }, [])

    if (isLoading) {
        return <h3>Loading...</h3>
    }

    return (
        <div className="topic-health-box">
            <div className="topic-title-box">
                <h3 className="topic-h3">Topics</h3>
            </div>
            <div className="topic-health">
                {availableTopics.map((topic) => {
                    return <div className="topic">
                        <HealthStatus isAlive={topic.status} />
                        <h4 className="topic-h3">{topic.topic}</h4>
                    </div>
                })}
            </div>
        </div >
    )
}