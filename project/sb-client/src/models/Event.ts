import { Line } from "./Line";

export interface Event {
    id: number,
    eventDate: Date,
    eventType: number,
    isFrozen: boolean,
    relatedGame: {
        id: number,
        name: string,
    },
    league: {
        id: number,
        name: string,
        country: string,
    },
    lines: Line[]
}