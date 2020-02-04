import axios from "axios";

export default class Api {

    static axiosInstance = axios.create({
        baseURL: process.env.REACT_APP_API_URL || "http://localhost:8090",
        headers: {'Content-Type' : 'application/json', "Authorization": "Bearer token"}
    });

    static getSpotifyUrl() {
        return this.axiosInstance.get("/spotify")
    }
}
