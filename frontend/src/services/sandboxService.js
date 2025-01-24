import axios from "axios";

class SandboxService {

    constructor() {
        this.apiClient = axios.create({
            baseURL: "http://localhost:8080",
            headers: {
                "Content-Type": "application/json",
                "Access-Control-Allow-Origin": "*",
            },
        });
    }

    async getAllSandboxes() {
        const response = await this.apiClient.get("/api/sandboxes");
        console.log(response)
        return response.data;
    }

    async deleteSandbox(id) {
        const response = await this.apiClient.delete(`/api/sandboxes/${id}`);
        return response.data;
    }

    async createSandbox(image_name, lifetime) {
        var data = {
            "image_name": image_name,
            "lifetime": lifetime
        };
        const response = await this.apiClient.post("/api/sandboxes", data);
        return response.data;
    }

}

export default new SandboxService();