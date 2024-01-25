import "./styles.css"


export const HealthStatus: React.FC<{ isAlive: boolean }> = ({ isAlive }) => {
    return <div className="health">
        <div className={isAlive ? 'status-alive' : 'status-dead'} />
    </div>
}