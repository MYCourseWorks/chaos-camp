import React from 'react';
import { Route, Redirect, RouteProps } from 'react-router-dom';
import { MainStore } from '../stores/mainStore';

interface Props extends RouteProps {
    mainStore: MainStore
}

export class PrivateRoute extends React.Component<Props> {
    public render(): React.ReactNode {
        const isAuth = this.props.mainStore.services.auth.IsAuthenticated
        if (!isAuth) {
            return (
                <Route
                    {...this.props}
                    component={() => <Redirect to="/login" />}
                    render={undefined}
                />
            )
        } else {
            return (<Route {...this.props} />)
        }
    }
}