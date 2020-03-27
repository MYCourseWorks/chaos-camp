import { User } from "../models/User";

export class Users {
    public async all(siteID: number): Promise<User[]> {
        const resp = await fetch(
            `http://localhost:8080/api/v1/users/all/${siteID}`,
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
            const body: User[] = await resp.json();
            return body;
        } catch(err) {
            throw err;
        }
    }

    public async delete(userID: number): Promise<void> {
        const resp = await fetch(
            `http://localhost:8080/api/v1/users/delete/${userID}`,
            {
                method: "Delete",
                headers: new Headers({
                    "Authorization": "Bearer " + localStorage.getItem("token"),
                }),
            },
        );

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`);
        }
    }

    public async create(editModel: User): Promise<number> {
        const resp = await fetch(
            `http://localhost:8080/api/v1/users/create`,
            {
                method: "Post",
                headers: new Headers({
                    "Authorization": "Bearer " + localStorage.getItem("token"),
                }),
                body: JSON.stringify(editModel),
            },
        );

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`);
        }

        try {
            const user: { id: number } = await resp.json();
            if (!user.id) {
                throw new Error("Invalid respnse")
            }
            return user.id
        } catch(err) {
            throw err;
        }
    }

    public async update(userID: number, username: string, password: string): Promise<void> {
        const body: { userID: number, username?: string, password?: string } = { userID: userID }
        if (username) body.username = username
        if (password) body.password = password

        const resp = await fetch(
            `http://localhost:8080/api/v1/users/update`,
            {
                method: "Put",
                headers: new Headers({
                    "Authorization": "Bearer " + localStorage.getItem("token"),
                }),
                body: JSON.stringify(body)
            },
        );

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`);
        }
    }
}