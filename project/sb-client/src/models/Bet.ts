export interface Bet {
    id: number,
    userID: number,
    userName: string,
    value: string,
    isPayed: boolean,
    betOnIndex: number,
    lineID: number,
    lineType: number,
    lineDescription: string,
    gameID: number,
    gameName: string,
}