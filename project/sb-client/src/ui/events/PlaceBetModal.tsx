import React from "react"
import bindThis from "../../common/attributes/bindThis"

interface Props {
    isOpen: boolean
    title: string
    onCancel: () => void
    onSubmit: (betValue: string) => void
}

interface State {
    placeBetValue: string
}

export class PlaceBetModal extends React.Component<Props, State> {

    private initState: State = {
        placeBetValue: ""
    }

    constructor(props: Props) {
        super(props)
        this.state = { ...this.initState }
    }

    public render(): React.ReactNode {
        const { isOpen, onCancel, title } = this.props

        return (
            <div className="modal" style={{ display: isOpen ? "block" : "none" }}>
                <div className="modal-background" />
                <div className="modal-card">
                    <header className="modal-card-head">
                        <p className="modal-card-title">Place Bet on {title}</p>
                        <button className="delete" onClick={onCancel} />
                    </header>
                    <section className="modal-card-body">
                        <div className="field">
                            <label className="label">Bet Value</label>
                            <div className="control has-icons-left">
                                <input
                                    type="number"
                                    className="input"
                                    value={this.state.placeBetValue}
                                    min="0"
                                    onChange={(e) =>
                                        this.setState({ placeBetValue: e.target.value })
                                    }
                                    required
                                />
                            </div>
                        </div>
                    </section>
                    <footer className="modal-card-foot">
                        <button className="button" onClick={this.submitClicked}>Ok</button>
                        <button className="button" onClick={onCancel}>Cancel</button>
                    </footer>
                </div>
            </div>
        )
    }

    @bindThis
    private submitClicked(): void {
        try {
            if (this.state.placeBetValue) {
                this.props.onSubmit(this.state.placeBetValue)
            } else {
                console.error("Place bet value is invalid!")
            }
        } catch (err) {
            console.error(err)
        } finally {
            this.setState({ ...this.initState })
        }
    }
}