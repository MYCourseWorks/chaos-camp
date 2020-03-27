import { League } from "../models/League";

export class Leagues {
    public async all(): Promise<League[]> {
        const resp = await fetch(
            "http://localhost:8081/api/v1/leagues/all",
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
            const body: League[] = await resp.json();
            return body;
        } catch(err) {
            throw err;
        }
    }
}