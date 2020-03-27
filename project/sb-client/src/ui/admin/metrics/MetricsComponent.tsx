import React from "react"
import { MainStore } from "../../../stores/mainStore"
import { Metric } from "../../../models/Metric"

interface Props {
    mainStore: MainStore
}

interface State {
    metrics: Metric[]
}

export class MetricsComponent extends React.Component<Props, State> {

    constructor(props: Props) {
        super(props)
        this.state = { metrics: [] }
    }

    public async componentDidMount(): Promise<void> {
        const metrics: Metric[] = await this.props.mainStore.services.metrics.all()
        this.setState({ metrics })
    }

    public render(): React.ReactNode {
        const metrics = this.state.metrics

        return (
            <div className="panel">
                <div
                    className="panel-heading"
                    style={{ textAlign: "center" }}
                >
                    Metrics
                </div>
                <div className="panel-block has-text-white has-background-primary has-text-weight-bold">
                    <div className="column is-3">Timestamp</div>
                    <div className="column is-3">Method</div>
                    <div className="column is-3">URL</div>
                    <div className="column is-3">Response</div>
                </div>
                {
                    this.state.metrics.map((m: Metric) =>
                        <div className="panel-block" key={m.id}>
                            <div className="column is-3">{new Date(m.timestamp).toLocaleString()}</div>
                            <div className="column is-3">{m["http-method"]}</div>
                            <div className="column is-3">{m.url}</div>
                            <div className="column is-3">{m["resp-formatted"]}</div>
                        </div>
                    )
                }
            </div>
        )
    }
}