import { Auth } from '../services/auth'
import { Leagues } from '../services/leagues'
import { Users } from '../services/users'
import { Events } from '../services/events'
import { Bets } from '../services/bets'
import { Metrics } from '../services/metrics'

export class MainStore {
    public services: {
        auth: Auth,
        leagues: Leagues,
        users: Users,
        games: Events,
        bets: Bets,
        metrics: Metrics,
    }

    constructor() {
        this.services = {
            auth: new Auth(),
            leagues: new Leagues(),
            users: new Users(),
            games: new Events(),
            bets: new Bets(),
            metrics: new Metrics(),
        }
    }
}