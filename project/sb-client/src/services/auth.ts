export class Auth {
    public get IsAuthenticated(): boolean {
        return !!localStorage.getItem("token")
    }

    public get IsSiteUser(): boolean {
        const rolesStr = localStorage.getItem("roles")
        if (rolesStr) {
            return parseInt(rolesStr) < 4
        } else {
            throw new Error("No roles in local storage")
        }
    }

    public get IsAdmin(): boolean {
        const rolesStr = localStorage.getItem("roles")
        if (rolesStr) {
            return parseInt(rolesStr) >= 8
        } else {
            throw new Error("No roles in local storage")
        }
    }

    public get IsOperator(): boolean {
        const rolesStr = localStorage.getItem("roles")
        if (rolesStr) {
            return parseInt(rolesStr) >= 4
        } else {
            throw new Error("No roles in local storage")
        }
    }

    public async login(name: string, password: string, siteID: number): Promise<void> {
        const resp = await fetch("http://localhost:8080/api/v1/users/login", {
            method: "Post",
            body: JSON.stringify({ name, password, siteID }),
        });

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`)
        }

        try {
            const body: { token: string, roles: number, id: number } = await resp.json()
            if (!body.token) {
                throw new Error("No token in response!")
            }

            localStorage.setItem("user_id", body.id.toString())
            localStorage.setItem("user_name", name)
            localStorage.setItem("roles", body.roles.toString())
            localStorage.setItem("token", body.token)
        } catch(err) {
            throw err
        }
    }

    public async register(name: string, password: string, siteID: number): Promise<void> {
        const resp = await fetch("http://localhost:8080/api/v1/users/register", {
            method: "Post",
            body: JSON.stringify({ name, password, siteID }),
        });

        if (!resp.ok) {
            throw new Error(`${resp.status} ${resp.statusText}`)
        }

        try {
            const body: { id: number, roles: number } = await resp.json()
            if (!body.id || !body.roles) {
                throw new Error("Register failed!")
            }

            localStorage.setItem("user_id", body.id.toString())
            localStorage.setItem("roles", body.roles.toString())
        } catch(err) {
            throw err
        }
    }

    public logout(): void {
        localStorage.clear()
        window.location.reload(true)
    }
}