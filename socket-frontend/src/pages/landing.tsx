import { useEffect, useState } from "react"
import { HubDetail } from "../types/connection_detail";
import HubDetails from "../components/hub_detail/hub_detail";
import AppBar from "../components/appbar/appbar";
import { getHubDetails } from "../service/api_service";
import TopicHealth from "../components/topic_health/topic_health";

export default function Landing() {

    const [isLoading, setLoading] = useState(true)
    const [hubDetail, setDetail] = useState<HubDetail>({
        active_connections: 2,
        connections: []
    });

    useEffect(() => {
        setTimeout(async () => {
            setDetail(await getHubDetails())
            setLoading(false)
        }, 1000);
    }, [])

    if (isLoading) {
        return <h2>Loading..</h2>
    }

    return <div>
        <AppBar/>
        <TopicHealth/>
        {HubDetails(hubDetail)}
    </div>;
}