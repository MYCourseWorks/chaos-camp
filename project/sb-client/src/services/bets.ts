import { Bet } from "../models/Bet";

export class Bets {
    public async allForUser(userID: number): Promise<Bet[]> {
        const resp = await fetch(
            `http://localhost:8082/api/v1/bets/${userID}`,
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
            const body: Bet[] = await resp.json();
            return body;
        } catch(err) {
            throw err;
        }
    }

    public async all(): Promise<Bet[]> {
        const resp = await fetch(
            `http://localhost:8082/api/v1/bets/all`,
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
            const body: Bet[] = await resp.json();
            return body;
        } catch(err) {
            throw err;
        }
    }

    public async cancel(userID: number, betID: number): Promise<void> {
        const resp = await fetch(
            `http://localhost:8082/api/v1/bets/cancel`,
            {
                method: "Delete",
                headers: new Headers({
                    "Authorization": "Bearer " + localStorage.getItem("token"),
                }),
                body: JSON.stringify({ betID, userID })
            },
        );

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`);
        }
    }

    public async payout(betID: number): Promise<void> {
        const resp = await fetch(
            `http://localhost:8082/api/v1/bets/payout`,
            {
                method: "Post",
                headers: new Headers({
                    "Authorization": "Bearer " + localStorage.getItem("token"),
                }),
                body: JSON.stringify({ betID })
            },
        );

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`);
        }
    }

    public async place(
        userID: number,
        lineID: number,
        betOnIndex: number,
        value: string,
    ): Promise<void> {

        const resp = await fetch(
            `http://localhost:8082/api/v1/bets/place`,
            {
                method: "Post",
                headers: new Headers({
                    "Authorization": "Bearer " + localStorage.getItem("token"),
                }),
                body: JSON.stringify({ userID, lineID, value, betOnIndex })
            },
        );

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`);
        }
    }
}