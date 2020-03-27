import React from "react"
import { Line } from "../../models/Line"
import "./Line.css"

interface Props {
    team1: string
    team2: string
    line: Line
    isFrozen: boolean
    placeBetClicked(lineIndex: number): void
}

interface State {
}

export class EventLine extends React.Component<Props, State> {

    public render(): React.ReactNode {
        const { line, team1, team2, placeBetClicked, isFrozen } = this.props
        const classForOdd = "has-text-right has-background-white-ter line-box-content"
        const classForType = "has-text-left has-background-white-ter line-box-content"

        return (
            <div className="columns line-container">
                <div className="column line-box">
                    <div
                        className={classForType}
                        onClick={() => {
                            if (!isFrozen) { placeBetClicked(0) }
                        }}
                    >
                        {team1}
                    </div>
                </div>
                <div className="column line-box">
                    <div className={classForOdd}>
                        {parseFloat(line.odds[0]).toFixed(2)}
                    </div>
                </div>
                <div className="column line-box">
                    <div
                        className={classForType}
                        onClick={() => {
                            if (!isFrozen) { placeBetClicked(1) }
                        }}
                    >
                        Draw
                    </div>
                </div>
                <div className="column line-box">
                    <div className={classForOdd}>
                        {parseFloat(line.odds[1]).toFixed(2)}
                    </div>
                </div>
                <div className="column line-box">
                    <div
                        className={classForType}
                        onClick={() => {
                            if (!isFrozen) { placeBetClicked(2) }
                        }}
                    >
                        {team2}
                    </div>
                </div>
                <div className="column line-box">
                    <div className={classForOdd}>
                        {parseFloat(line.odds[2]).toFixed(2)}
                    </div>
                </div>
            </div>
        )
    }
}