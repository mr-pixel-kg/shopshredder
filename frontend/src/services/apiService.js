import axios from "axios";

class ApiService {

    constructor() {
        this.apiClient = axios.create({
            baseURL: "http://localhost:8080",
            headers: {
                "Content-Type": "application/json",
                "Access-Control-Allow-Origin": "*",
            },
        });
        this.authCredentials = null;
    }

    async login(username, password) {
        const response = await this.apiClient.get("/api/auth", {
            auth: {
                username: username,
                password: password
            }
        });
        console.log(response)

        if (response.data.loggedIn === "true") {
            this.authCredentials = { username, password };
            return true;
        }

        return false;
    }

    logout() {
        this.authCredentials = null;
    }

    isLoggedIn() {
        return this.authCredentials !== null;
    }

}

export default new ApiService();