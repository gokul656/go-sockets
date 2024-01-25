import axios from "axios";
import { HubDetail, TopicData } from "../types/connection_detail";
import { AVAILABLE_TOPIC_ROUTE, HUB_DETAIL_ROUTE } from "./routes";

const BASE_URL = "http://192.168.201.138:8080/api"

export async function getHubDetails(): Promise<HubDetail> {
    try {
        const response = await axios.get(`${BASE_URL}${HUB_DETAIL_ROUTE}`);
        return response.data;
    } catch (e) {
        console.log("ERROR", e)
        const defaultValue: HubDetail = {
            active_connections: 0,
            connections: []
        };
        return Promise.resolve(defaultValue)
    }
}


export async function getAvailableTopics(): Promise<TopicData[]> {
    try {
        const response = await axios.get(`${BASE_URL}${AVAILABLE_TOPIC_ROUTE}`);
        return response.data;
    } catch (e) {
        console.log("ERROR", e)
        return Promise.resolve([])
    }
}