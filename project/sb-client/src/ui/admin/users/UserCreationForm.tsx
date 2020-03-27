import React from "react"
import { User } from "../../../models/User"

interface Props {
    editModel: User
    onChangeUsername(value: string): void
    onChangePassword(value: string): void
    onChangeRole(value: number): void
}

export class UserCreationForm extends React.Component<Props> {

    public render(): React.ReactNode {
        const editModel = this.props.editModel

        return (
            <div>
                <div className="field">
                    <label className="label">Name</label>
                    <div className="control">
                        <input
                            className="input"
                            type="text"
                            value={editModel.name}
                            onChange={(event) => {
                                this.props.onChangeUsername(event.target.value)
                            }}
                        />
                    </div>
                </div>
                <div className="field">
                    <label className="label">Password</label>
                    <div className="control">
                        <input
                            className="input"
                            type="password"
                            value={editModel.password}
                            onChange={(event) => {
                                this.props.onChangePassword(event.target.value)
                            }}
                        />
                    </div>
                </div>
                <div className="field">
                    <label className="label">Roles</label>
                    <div className="control">
                        <div className="select">
                            <select
                                value={editModel.roles}
                                onChange={(event) => {
                                    this.props.onChangeRole(parseInt(event.target.value))
                                }}
                            >
                                <option value="2">User</option>
                                <option value="4">Operator</option>
                                <option value="8">Admin</option>
                            </select>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}