import { Metric } from "../models/Metric";

export class Metrics {
    public async all(): Promise<Metric[]> {
        const resp = await fetch(
            "http://localhost:8083/api/v1/metrics/all",
            {
                method: "Get",
                headers: new Headers({
                    "Authorization": "Bearer " + localStorage.getItem("token"),
                }),
            },
        );

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`);
        }

        try {
            const body: Metric[] = await resp.json();
            return body;
        } catch(err) {
            throw err;
        }
    }
}