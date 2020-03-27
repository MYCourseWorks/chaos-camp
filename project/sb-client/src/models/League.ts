import { Sport } from "./Sport";

export interface League {
    id: number,
    name: string,
    sport: Sport,
    country: string,
}