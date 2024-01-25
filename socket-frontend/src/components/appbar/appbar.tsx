import "./styles.css"
import { HealthStatus } from "../status/status";

export default function AppBar() {
    return <div className="socket-overview">
        <h3>Dashboard</h3>
        <HealthStatus isAlive={true}/>
    </div>
}