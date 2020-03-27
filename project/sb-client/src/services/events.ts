import { Event } from "../models/Event";

export class Events {
    public async allPerLeague(leagueID: number): Promise<Event[]> {
        const resp = await fetch(
            `http://localhost:8081/api/v1/games/all?leagueID=${leagueID}`,
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
            const body: Event[] = await resp.json();
            return body;
        } catch(err) {
            throw err;
        }
    }

    public async toggleFreeze(eventID: number): Promise<void> {
        const resp = await fetch(
            `http://localhost:8081/api/v1/games/freeze/${eventID}`,
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
    }
}