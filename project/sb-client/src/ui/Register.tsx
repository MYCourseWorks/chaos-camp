import React from "react"
import { Link } from "react-router-dom"
import bindThis from "../common/attributes/bindThis"
import { MainStore } from "../stores/mainStore";

interface State {
    name: string,
    password: string,
}

interface Props {
    mainStore: MainStore
}

export class Register extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = { name: "", password: "" }
    }

    @bindThis
    public changeName(event: React.ChangeEvent<HTMLInputElement>) {
        this.setState({name: event.target.value})
    }

    @bindThis
    public changePassword(event: React.ChangeEvent<HTMLInputElement>) {
        this.setState({password: event.target.value})
    }

    public render(): React.ReactNode {
        return (
            <div className="container">
                <div className="box">
                    <div className="title">
                        <div className="title is-parent is-vertical">
                            <div className="tile is-child notification is-warning">
                                <p className="title">Register</p>
                            </div>
                        </div>
                    </div>

                    <div className="field">
                        <label className="label">Name</label>
                        <div className="control has-icons-left">
                            <input
                                type="name"
                                className="input"
                                value={this.state.name}
                                onChange={this.changeName}
                                required
                            />
                            <span className="icon is-small is-left">
                                <i className="fas fa-user" />
                            </span>
                        </div>
                    </div>
                    <div className="field">
                        <label className="label">Password</label>
                        <div className="control has-icons-left">
                            <input
                                type="password"
                                className="input"
                                value={this.state.password}
                                onChange={this.changePassword}
                                required
                            />
                            <span className="icon is-small is-left">
                                <i className="fa fa-lock" />
                            </span>
                        </div>
                    </div>
                    <div className="field">
                        <button className="button is-success" onClick={this.register}>Submit</button>
                        <Link to="/login" className="button is-cancel" style={{marginLeft:"5px"}}>Login</Link>
                    </div>
                </div>
            </div>
        )
    }

    @bindThis
    private async register(): Promise<void> {
        if (this.state.name && this.state.password) {
            try {
                const auth = this.props.mainStore.services.auth
                await auth.register(this.state.name, this.state.password, 1)
                window.location.replace("/login")
            } catch (err) {
                console.error(err)
            }
        }
    }
}
