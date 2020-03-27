import React from "react"
import "./Navigation.css"
import { Link } from "react-router-dom"
import { MainStore } from "../../stores/mainStore"


interface Props {
    mainStore: MainStore
}

export class Navigation extends React.Component<Props> {

    public render(): React.ReactNode {
        const auth = this.props.mainStore.services.auth

        return (
            <nav className="navbar app-navigation" role="navigation" aria-label="main navigation">
                <div className="navbar-brand">
                    <Link className="navbar-item" to="/">
                        <img src="https://bulma.io/images/bulma-logo.png" width={112} height={28} alt="logo" />
                    </Link>
                    <button className="navbar-burger burger" data-target="navbarBasicExample">
                        <span aria-hidden="true" />
                        <span aria-hidden="true" />
                        <span aria-hidden="true" />
                    </button>
                </div>
                <div id="navbarBasicExample" className="navbar-menu">
                    <div className="navbar-start">
                        <Link className="navbar-item" to="/">Home</Link>
                        <Link className="navbar-item" to="/profile">Profile</Link>
                        <Link className="navbar-item" to="/bets">Bets</Link>
                        { auth.IsAdmin && <Link className="navbar-item" to="/users/all">Users</Link> }
                        { auth.IsAdmin && <Link className="navbar-item" to="/metrics/all">Metrics</Link> }
                    </div>
                    <div className="navbar-end">
                        <div className="navbar-item">
                            <div className="buttons">
                                <button className="button is-light" onClick={auth.logout}>Logout</button>
                            </div>
                        </div>
                    </div>
                </div>
            </nav>
        )
    }
}