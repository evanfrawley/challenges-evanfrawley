// React stuff
import React, {Component} from 'react';
import {Route, Switch, withRouter} from 'react-router'
import {Navbar, NavItem} from 'react-materialize';

// Custom components
import Login from './components/Login';
import SignUp from './components/SignUp';
import Home from './components/Home';
import Settings from './components/Settings';
import Messaging from './components/messaging/Messaging';
import Search from "./components/Search";
import NoMatch from "./components/NoMatch";

// Routing
import {PrivateRoute} from './routing/RoutingHelpers';

// Helpers
import 'whatwg-fetch';
import * as AuthService from './services/AuthAPIService';
import * as Helpers from './services/APIHelpers';

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
      this.loadCurrentUserData();
    }
  }

  handleLogin = (credentials) => {
    AuthService.signInUser(credentials).then((response) => {
      if (response !== '') {
        this.setState({loggedIn: true, authToken: response})
      }
      this.props.history.push('/');
    }).then(() => {
      this.loadCurrentUserData();
    }).catch((error) => {
        console.error('error when logging in user', error);
        AuthService.removeLocalStorage();
        this.props.history.push('/');
    })
  };

  loadCurrentUserData = () => {
    AuthService.getUser().then((response) => {
      this.setState({user: response});
    }).catch((error) => {
        console.error('error when logging in user', error);
        this.props.history.push('/');
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
        this.setState({loggedIn: false});
        AuthService.removeLocalStorage();
      });
  };

  handleNavigateToSettings = () => {
    this.props.history.push('/settings')
  };

  render() {
    return (
      <div className="App">
        <Navbar brand='Messages' right>
          <NavItem href='/messaging'>Messaging</NavItem>
          <NavItem href='/login'>Login</NavItem>
        </Navbar>
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
          <PrivateRoute
            path='/messaging/:channelID'
            redirectTo={'/login'}
            component={Messaging}
            authed={this.state.loggedIn}
            user={this.state.user}
            handleSettingsUpdate={this.handleSettingsUpdate}
          />
          <PrivateRoute
            path='/messaging'
            redirectTo={'/login'}
            component={Messaging}
            authed={this.state.loggedIn}
            user={this.state.user}
            handleSettingsUpdate={this.handleSettingsUpdate}
          />
          <Route component={NoMatch}/>
        </Switch>
      </div>
    );
  }
}

export default withRouter(App);
