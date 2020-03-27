import React from "react"
import { MainStore } from "../../stores/mainStore"
import bindThis from "../../common/attributes/bindThis"


interface Props {
    mainStore: MainStore
}

interface State {
    name: string,
    password: string,
}

export class UserProfile extends React.Component<Props, State> {

    constructor(props: Props) {
        super(props);
        this.state = { name: "", password: "" }
    }

    public render(): React.ReactNode {
        return (
            <div className="container">
                <div className="box">
                    <div className="title">
                        <div className="title is-parent is-vertical">
                            <div className="tile is-child notification is-primary">
                                <p className="title">Change Profile</p>
                            </div>
                        </div>
                    </div>

                    <div className="field">
                        <label className="label">Change Name</label>
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
                        <label className="label">Change Password</label>
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
                        <button className="button is-success" onClick={this.update}>Submit</button>
                    </div>
                </div>
            </div>
        )
    }

    @bindThis
    public changeName(event: React.ChangeEvent<HTMLInputElement>) {
        this.setState({name: event.target.value})
    }

    @bindThis
    public changePassword(event: React.ChangeEvent<HTMLInputElement>) {
        this.setState({password: event.target.value})
    }


    @bindThis
    private async update(): Promise<void> {
        if (this.state.name || this.state.password) {
            try {
                const users = this.props.mainStore.services.users
                const userID = localStorage.getItem("user_id")
                if (!userID) throw new Error("No userID in localStorage")
                await users.update(parseInt(userID), this.state.name, this.state.password)
                window.location.reload(true)
            } catch(err) {
                console.error(err)
            }
        }
    }
}