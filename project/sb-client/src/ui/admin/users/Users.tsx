import React from "react"
import { MainStore } from "../../../stores/mainStore"
import { User } from "../../../models/User"
import { ConfirmModal } from "../../modal/ConfirmModal"
import bindThis from "../../../common/attributes/bindThis"
import { UserCreationForm } from "./UserCreationForm"

interface Props {
    mainStore: MainStore
}

interface State {
    users: User[]
    isUserDelModalOpened: boolean
    isUserCreateModalOpened: boolean
    editModel?: User
}

export class Users extends React.Component<Props, State> {
    private initEditModel: User = { id: -1, name: "", roles: 2, siteID: 1, password: "" }

    constructor(props: Props) {
        super(props)
        this.state = {
            users: [],
            isUserDelModalOpened: false,
            isUserCreateModalOpened: false,
            editModel: { ...this.initEditModel },
        }
    }

    public async componentDidMount(): Promise<void> {
        const users: User[] = await this.props.mainStore.services.users.all(1)
        this.setState({ users })
    }

    public render(): React.ReactNode {
        return (
            <div className="panel">
                <div className="panel-heading" style={{ textAlign: "center" }}>Users</div>
                {
                    this.state.users.map((u: User) =>
                        <div className="panel-block" key={u.id}>
                            <div className="column is-11">
                                <span style={{ marginRight: "10px" }}>
                                    {this.pickUserIcon(u)}
                                </span>
                                {u.name}
                            </div>
                            <div className="column is-1" style={{ textAlign: "right", cursor: "pointer" }}>
                                <span onClick={this.openDelModal(u)}>
                                    <i className="fas fa-trash" aria-hidden="true" />
                                </span>
                            </div>
                        </div>
                    )
                }
                <div className="panel-block">
                    <button className="button is-link is-outlined is-fullwidth"
                        onClick={this.openCreateModal()}
                    >
                        Create New User
                    </button>
                </div>

                <ConfirmModal
                    isOpen={this.state.isUserDelModalOpened}
                    title={"Delete User"}
                    onCancel={this.onDeleteUserCancel}
                    onSubmit={this.onDeleteUserSubmit}
                >
                    <h1>Are you sure you want to delete this user ?</h1>
                </ConfirmModal>

                <ConfirmModal
                    isOpen={this.state.isUserCreateModalOpened}
                    title={"Delete User"}
                    onCancel={this.onUserCreateCancel}
                    onSubmit={this.onUserCreateSubmit}
                >
                    <UserCreationForm
                        editModel={this.state.editModel as User}
                        onChangePassword={(value: string) => {
                            this.setState((prev) => {
                                if (!prev.editModel) { throw new Error("invalid edit model") }
                                prev.editModel.password = value
                                return { ...prev }
                            })
                        }}
                        onChangeRole={(value: number) => {
                            this.setState((prev) => {
                                if (!prev.editModel) { throw new Error("invalid edit model") }
                                prev.editModel.roles = value
                                return { ...prev }
                            })
                        }}
                        onChangeUsername={(value: string) => {
                            this.setState((prev) => {
                                if (!prev.editModel) { throw new Error("invalid edit model") }
                                prev.editModel.name = value
                                return { ...prev }
                            })
                        }}
                    />
                </ConfirmModal>
            </div>
        )
    }

    @bindThis
    public onDeleteUserCancel() {
        this.setState({ editModel: { ...this.initEditModel }, isUserDelModalOpened: false })
    }

    @bindThis
    public onUserCreateCancel() {
        this.setState({ editModel: { ...this.initEditModel }, isUserCreateModalOpened: false })
    }

    @bindThis
    public async onDeleteUserSubmit(): Promise<void> {
        if (this.state.editModel) {
            const id = this.state.editModel?.id
            try {
                await this.props.mainStore.services.users.delete(id as number)
                window.location.reload(true)
            } catch (err) {
                console.error(err)
            }
        } else {
            console.error("Deletion failed")
        }

        this.setState({ isUserDelModalOpened: false, editModel: { ...this.initEditModel } })
    }

    @bindThis
    public async onUserCreateSubmit(): Promise<void> {
        if (this.state.editModel) {
            const editModel = this.state.editModel
            try {
                const userId = await this.props.mainStore.services.users.create(editModel)
                editModel.id = userId
                this.setState({ users: this.state.users.concat(editModel) })
            } catch (err) {
                console.error(err)
            }
        } else {
            console.error("Creation failed")
        }

        this.setState({ isUserCreateModalOpened: false, editModel: { ...this.initEditModel } })
    }

    @bindThis
    public openDelModal(u: User): ((event: React.MouseEvent<HTMLSpanElement, MouseEvent>) => void) | undefined {
        return (event: React.MouseEvent<HTMLSpanElement, MouseEvent>) => {
            this.setState({ isUserDelModalOpened: true, editModel: u })
        }
    }

    @bindThis
    public openCreateModal(): ((event: React.MouseEvent<HTMLSpanElement, MouseEvent>) => void) | undefined {
        return (event: React.MouseEvent<HTMLSpanElement, MouseEvent>) => {
            this.setState({ isUserCreateModalOpened: true, editModel: { ...this.initEditModel } })
        }
    }

    public pickUserIcon(u: User): React.ReactNode {
        if (u.roles >= 8) return <i className="fas fa-user-shield"></i>
        else if (u.roles >= 4) return <i className="fas fa-users-cog"></i>
        else return <i className="fas fa-user"></i>
    }
}