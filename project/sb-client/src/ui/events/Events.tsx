import React from "react"
import { MainStore } from "../../stores/mainStore"
import { Event } from "../../models/Event"
import { EventLine } from "./Line"
import { Line } from "../../models/Line"
import { PlaceBetModal } from "./PlaceBetModal"
import bindThis from "../../common/attributes/bindThis"

interface Props {
    leagueID: number
    mainStore: MainStore
}

interface State {
    events: Event[]
    isPlaceBetOpened: boolean
    placeBetIndex?: number
    placeBetLineID?: number
    placeBetTitle: string
}

export class Events extends React.Component<Props, State> {

    private initState: Partial<State> = {
        isPlaceBetOpened: false,
        placeBetIndex: undefined,
        placeBetLineID: undefined,
        placeBetTitle: "",
    }

    constructor(props: Props) {
        super(props)
        this.state = { events: [], ...this.initState } as State
    }

    public async componentWillMount(): Promise<void> {
        await this.loadGamse(this.props)
    }

    public async componentWillReceiveProps(nextProps: Props): Promise<void> {
        await this.loadGamse(nextProps)
    }

    public async loadGamse(props: Props): Promise<void> {
        const leagueID = props.leagueID
        const event = await props.mainStore.services.games.allPerLeague(leagueID)
        this.setState({ events: event })
    }

    public render(): React.ReactNode {
        return (
            <div>
                {
                    this.state.events.length > 0 &&
                    <h1 className="title">{this.state.events[0].league.name}</h1>
                }
                <ul>
                    {
                        this.state.events &&
                        this.state.events.length > 0 &&
                        this.state.events.map((e: Event) =>
                            <li key={e.id} className="box" style={{
                                padding: 0,
                                opacity: e.isFrozen ? 0.3 : 1,
                            }}>
                                <div style={{ padding: "1%" }}>
                                    <div className="columns">
                                        <div className="column has-background-primary has-text-white">
                                            {e.relatedGame.name}
                                        </div>
                                        {
                                            this.props.mainStore.services.auth.IsOperator &&
                                            <div
                                                className="column has-background-primary has-text-white has-text-right"
                                                onClick={() => this.freezeEvent(e.id)}
                                            >
                                                <i className="fas fa-snowflake"></i>
                                            </div>
                                        }
                                    </div>
                                    {e.lines && e.lines.map((l: Line) =>
                                        <EventLine
                                            key={l.id}
                                            line={l}
                                            team1={e.relatedGame.name.split(" vs ")[0].trim()}
                                            team2={e.relatedGame.name.split(" vs ")[1].trim()}
                                            isFrozen={e.isFrozen}
                                            placeBetClicked={(n: number) => {
                                                let title: string
                                                switch (n) {
                                                    case 0:
                                                        title = e.relatedGame.name.split(" vs ")[0].trim() + " to win"
                                                        break
                                                    case 1:
                                                        title = "Draw"
                                                        break
                                                    case 2:
                                                        title = e.relatedGame.name.split(" vs ")[1].trim() + " to win"
                                                        break
                                                    default:
                                                        throw new Error("Invalid bet index")
                                                }

                                                this.setState({
                                                    placeBetIndex: n,
                                                    placeBetLineID: l.id,
                                                    isPlaceBetOpened: true,
                                                    placeBetTitle: title,
                                                })
                                            }}
                                        />
                                    )}
                                </div>
                            </li>
                        )
                    }
                </ul>

                <PlaceBetModal
                    isOpen={this.state.isPlaceBetOpened}
                    onCancel={this.onPlaceBetCancel}
                    onSubmit={this.onPlaceBet}
                    title={this.state.placeBetTitle}
                />
            </div>
        )
    }

    @bindThis
    private onPlaceBetCancel(): void {
        this.setState({ ...this.initState } as State)
    }

    @bindThis
    private  async onPlaceBet(betValue: string): Promise<void> {
        try {
            const userID = localStorage.getItem("user_id") as string
            const { placeBetIndex, placeBetLineID } = this.state

            await this.props.mainStore.services.bets.place(
                parseInt(userID),
                placeBetLineID as number,
                placeBetIndex as number, betValue
            )
        } catch(err) {
            console.error(err)
        } finally {
            this.setState({ ...this.initState } as State)
        }
    }

    @bindThis
    private async freezeEvent(id: number): Promise<void> {
        try {
            await this.props.mainStore.services.games.toggleFreeze(id)
            window.location.reload(true)
        } catch (err) {
            console.error('Freeze Failed')
            console.error(err)
        }
    }
}