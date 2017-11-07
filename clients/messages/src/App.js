import React, {Component} from 'react';
import Login from './components/Login';
import SignUp from './components/SignUp';
import Home from './components/Home';
import Settings from './components/Settings';
import 'whatwg-fetch';
import {Route, Switch, withRouter} from 'react-router'
import * as AuthService from './services/AuthAPIService';
import {PrivateRoute} from './routing/RoutingHelpers';
import * as Helpers from './services/APIHelpers';
import Search from "./components/Search";

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {
            // loggedIn: false,
            authToken: '',
            user: {
                firstName: '',
                lastName: '',
            }
        }
    }

    componentWillMount() {
        let oneHour = 60 * 60 * 1000;
        let lastCreated = localStorage.getItem(Helpers.TOKEN_KEY_CREATED);
        if (new Date() - new Date(lastCreated) < oneHour) {
            this.setState({loggedIn: true});
            this.loadCurrentUserDate();
        }
    }

    handleLogin = (credentials) => {
        AuthService.signInUser(credentials).then((response) => {
            if (response !== '') {
                this.setState({loggedIn: true, authToken: response})
            }
            this.props.history.push('/');
        }).then(() => {
            this.loadCurrentUserDate();
        })
    };

    loadCurrentUserDate = () => {
        AuthService.getUser().then((response) => {
            this.setState({user: response});
        })
    };

    handleSettingsUpdate = (userUpdates) => {
        AuthService.updateUser(userUpdates)
            .then((response) => {
            })
    };

    handleSignOut = () => {
        AuthService.signOutUser()
            .then(() => {
                this.setState({loggedIn: false})
            });
    };

    handleNavigateToSettings = () => {
        this.props.history.push('/settings')
    };


    render() {
        console.log(this.state.loggedIn);
        return (
            <div className="App">
                <Switch>
                    <PrivateRoute
                        exact path='/'
                        redirectTo={'/login'}
                        component={Home}
                        authed={this.state.loggedIn}
                        user={this.state.user}
                        handleSignOut={this.handleSignOut}
                        handleNavigateToSettings={this.handleNavigateToSettings}
                    />
                    <Route
                        path='/login'
                        render={routerProps => <Login {...routerProps} handleLogin={this.handleLogin}/>}
                        handleLogin={this.handleLogin}
                    />
                    <Route
                        path='/signup'
                        component={SignUp}
                    />
                    <PrivateRoute
                        path='/settings'
                        redirectTo={'/login'}
                        component={Settings}
                        authed={this.state.loggedIn}
                        user={this.state.user}
                        handleSettingsUpdate={this.handleSettingsUpdate}
                    />
                    <PrivateRoute
                        path='/search'
                        redirectTo={'/login'}
                        component={Search}
                        authed={this.state.loggedIn}
                        user={this.state.user}
                        handleSettingsUpdate={this.handleSettingsUpdate}
                    />
                </Switch>
            </div>
        );
    }
}

export default withRouter(App);
