import React from 'react';
import 'bulma/css/bulma.css';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { PrivateRoute } from './common/PrivateRouter';
import { Login } from './ui/Login';
import { Register } from './ui/Register';
import "./App.css"
import { SportsMenu } from './ui/sports-menu/SportsMenu';
import { MainStore } from './stores/mainStore';
import { Navigation } from './ui/navigation/Navigation';
import { NotFound } from './ui/not-found-page/NotFound';
import { Users } from './ui/admin/users/Users';
import { Events } from './ui/events/Events';
import { BetsComponent } from './ui/bets/Bets';
import { UserProfile } from './ui/profile/UserProfile';
import { MetricsComponent } from './ui/admin/metrics/MetricsComponent';

interface RootProps {
    mainStore: MainStore
}

class Root extends React.Component<RootProps> {
    public render(): React.ReactNode {
        return (
            <div className="root-container">
                {/* Header */}
                <div className="columns">
                    <div className="column is-full">
                        <Navigation mainStore={this.props.mainStore} />
                    </div>
                </div>

                {/* Body */}
                <div className="columns">
                    <div className="column is-3">
                        <SportsMenu mainStore={this.props.mainStore} />
                    </div>
                    <div className="column is-8">
                        {this.props.children}
                    </div>
                </div>
            </div>
        )
    }
}

export class App extends React.Component {

    private mainStore: MainStore;

    constructor(props: {}) {
        super(props)
        this.mainStore = new MainStore()
    }

    public render(): React.ReactNode {
        return (
            <BrowserRouter>
                <Switch>
                    <Route path="/login" render={() => <Login mainStore={this.mainStore} />} />
                    <Route path="/register" render={() => <Register mainStore={this.mainStore} />} />
                    <PrivateRoute
                        exact path="/"
                        mainStore={this.mainStore}
                        render={(props) => (
                            <Root mainStore={this.mainStore}>
                                <div className="hero is-primary">
                                    <div className="hero-body has-text-centered">
                                        <h1 className="title">
                                            Welcome {localStorage.getItem("user_name")}
                                        </h1>
                                        <h1 className="subtitle">
                                            to Bulma Sports Betting
                                        </h1>
                                    </div>
                                </div>
                            </Root>
                        )} />
                    <PrivateRoute
                        exact path="/users/all"
                        mainStore={this.mainStore}
                        render={(props) => (
                            <Root mainStore={this.mainStore}>
                                <Users mainStore={this.mainStore} />
                            </Root>
                        )} />
                    <PrivateRoute
                        exact path="/leagues/:leagueID"
                        mainStore={this.mainStore}
                        render={(props) => (
                            <Root mainStore={this.mainStore}>
                                <Events
                                    mainStore={this.mainStore}
                                    leagueID={parseInt(props.match.params["leagueID"])}
                                />
                            </Root>
                        )} />
                    <PrivateRoute
                        exact path="/bets"
                        mainStore={this.mainStore}
                        render={(props) => (
                            <Root mainStore={this.mainStore}>
                                <BetsComponent
                                    mainStore={this.mainStore}
                                />
                            </Root>
                        )} />
                    <PrivateRoute
                        exact path="/profile"
                        mainStore={this.mainStore}
                        render={(props) => (
                            <Root mainStore={this.mainStore}>
                                <UserProfile mainStore={this.mainStore} />
                            </Root>
                        )} />
                    <PrivateRoute
                        exact path="/metrics/all"
                        mainStore={this.mainStore}
                        render={(props) => (
                            <Root mainStore={this.mainStore}>
                                <MetricsComponent mainStore={this.mainStore} />
                            </Root>
                        )} />
                    <Route component={() => <NotFound />} />
                </Switch>
            </BrowserRouter>
        )
    }
}

export default App;
