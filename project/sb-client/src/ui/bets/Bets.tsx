import React from "react"
import { MainStore } from "../../stores/mainStore"
import { Bet } from "../../models/Bet"
import { ConfirmModal } from "../modal/ConfirmModal"
import bindThis from "../../common/attributes/bindThis"

interface Props {
    mainStore: MainStore
}

interface State {
    bets: Bet[]
    isPayoutBetOpened: boolean,
    isCancelBetOpened: boolean,
    betIDToDelete?: number,
    userIDToDelete?: number,
    betIdToPayout?: number,
}

export class BetsComponent extends React.Component<Props, State> {

    private initState: Partial<State> = {
        isCancelBetOpened: false,
        isPayoutBetOpened: false,
        betIDToDelete: undefined,
        userIDToDelete: undefined,
        betIdToPayout: undefined,
    }

    constructor(props: Props) {
        super(props)
        this.state = { bets: [], ...this.initState as State }
    }

    public async componentDidMount(): Promise<void> {
        const userID = localStorage.getItem("user_id")
        const roles = localStorage.getItem("roles")
        if (!userID) throw new Error("No userID in localStorage")
        if (!roles) throw new Error("No roles in localStorage")

        const service = this.props.mainStore.services.bets
        const rolesAsInt = parseInt(roles)
        let bets: Bet[]
        if (rolesAsInt <= 2) {
            bets = await service.allForUser(parseInt(userID))
        } else {
            bets = await service.all()
        }

        this.setState({ bets })
    }

    public render(): React.ReactNode {
        const auth = this.props.mainStore.services.auth

        return (
            <div className="panel">
                <div
                    className="panel-heading"
                    style={{ textAlign: "center" }}
                >
                    Bets
                </div>
                <div className="panel-block has-text-white has-background-primary has-text-weight-bold">
                    <div className="column is-3">GAME</div>
                    <div className="column is-2">OUTCOME</div>
                    <div className="column is-2">VALUE</div>
                    <div className="column is-2">USER</div>
                    <div className="column is-3">ACTIONS</div>
                </div>
                {
                    this.state.bets.map((b: Bet) =>
                        <div className="panel-block" key={b.id}>
                            <div className="column is-3">
                                {b.gameName}
                            </div>
                            <div className="column is-2">
                                {b.betOnIndex === 0 && 'Home to win'}
                                {b.betOnIndex === 1 && 'Draw'}
                                {b.betOnIndex === 2 && 'Away to win'}
                            </div>
                            <div className="column is-2">
                                Bet: {b.value}$
                            </div>
                            <div className="column is-2">
                                {b.userName}
                            </div>
                            {
                                (!b.isPayed)
                                    ? (
                                        <div className="column is-3" style={{ textAlign: "right" }}>
                                            {
                                                ((auth.IsAdmin || auth.IsOperator) && !b.isPayed) &&
                                                (
                                                    <button
                                                        className="button is-info"
                                                        style={{ marginRight: "5px" }}
                                                        onClick={() => this.payoutBetOpen(b.id)}
                                                    >
                                                        Payout
                                                    </button>
                                                )
                                            }
                                            <button
                                                className="button is-danger"
                                                onClick={() => this.cancelBetOpen(b.userID, b.id)}
                                            >
                                                Cancel
                                            </button>
                                        </div>
                                    ) :
                                    (
                                        <div className="column is-3" style={{ textAlign: "right" }}>
                                            Bet Was Payed
                                        </div>
                                    )
                            }
                        </div>
                    )
                }

                <ConfirmModal
                    isOpen={this.state.isPayoutBetOpened}
                    title={"Payout Bet"}
                    onCancel={() => { this.setState({ ...this.initState as State }) }}
                    onSubmit={this.onBetPayout}
                >
                    <h1> Are you sure you want to payout this bet? </h1>
                </ConfirmModal>

                <ConfirmModal
                    isOpen={this.state.isCancelBetOpened}
                    title={"Cancel Bet"}
                    onCancel={() => { this.setState({ ...this.initState as State }) }}
                    onSubmit={this.onCancelBet}
                >
                    <h1>Are you sure you want to cancel this bet ?</h1>
                </ConfirmModal>
            </div>
        )
    }

    @bindThis
    private payoutBetOpen(betID: number): void {
        this.setState({
            isPayoutBetOpened: true,
            betIdToPayout: betID,
        })
    }

    @bindThis
    private async onBetPayout(): Promise<void> {
        try {
            const { betIdToPayout } = this.state
            await this.props.mainStore.services.bets.payout(betIdToPayout as number)
            window.location.reload(true)
        } catch (err) {
            console.error(err)
        }

        this.setState({ ...this.initState as State })
    }

    @bindThis
    private cancelBetOpen(userID: number, betID: number): void {
        this.setState({
            isCancelBetOpened: true,
            userIDToDelete: userID,
            betIDToDelete: betID
        })
    }

    @bindThis
    private async onCancelBet(): Promise<void> {
        try {
            const { userIDToDelete, betIDToDelete } = this.state
            await this.props.mainStore.services.bets
                .cancel(userIDToDelete as number, betIDToDelete as number)
            window.location.reload(true)
        } catch (err) {
            console.error(err)
        }

        this.setState({ ...this.initState as State })
    }
}