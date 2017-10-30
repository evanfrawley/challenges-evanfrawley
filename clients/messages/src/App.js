import React, {Component} from 'react';
import Login from './components/Login';
import SignUp from './components/SignUp';
import 'whatwg-fetch';
import { Switch, Route } from 'react-router'

import * as GoApiService from './services/GoApiService.js';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {
            summary: null,
        }
    }

    handleOnSubmit(e) {
        e.preventDefault();
        console.log('event target', e.target.url.value);
        GoApiService.getSummaryResourcePromise(e.target.url.value).then((json) => {
            this.setState({summary: json})
        })
    }

    render() {
        return (
            <div className="App">
                <Switch>
                    <Route exact path='/' component={Home}/>
                    <Route path='/login' component={Login}/>
                    <Route path='/signup' component={SignUp}/>
                </Switch>
            </div>
        );
    }
}

const Home = () => (
  <div className="Home">
      <Login />
  </div>
);


export default App;
