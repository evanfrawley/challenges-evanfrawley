import React from 'react';
import {Button, Input, Row} from 'react-materialize';
import {withRouter} from 'react-router';

import {signUpNewUser} from '../services/AuthAPIService';

class SignUp extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            newUser: {},
        };
    }

    _handleSubmit = (e) => {
        e.preventDefault();
        signUpNewUser(this.state.newUser)
            .then(() => {
                this.props.history.push('/login')
            })
    };

    _handleChange = (e) => {
        let newStateNewUser = this.state.newUser;
        newStateNewUser[e.target.name] = e.target.value;
        this.setState({newUser: newStateNewUser})
    };

    render() {
        return (
            <div className={"SignUpContainer"}>
                <Row>
                    <Input onChange={this._handleChange} name={"firstname"} s={6} label="First Name"/>
                    <Input onChange={this._handleChange} name={"lastname"} s={6} label="Last Name"/>
                    <Input onChange={this._handleChange} name={"username"} label="Username" s={12}/>
                    <Input onChange={this._handleChange} name={"password"} type="password" label="Password" s={12}/>
                    <Input onChange={this._handleChange} name={"passwordconf"} type="password" label="Confirm Password" s={12}/>
                    <Input onChange={this._handleChange} name={"email"} type="email" label="Email" s={12}/>
                    <Button onClick={this._handleSubmit}>Submit!</Button>
                </Row>
                <div>
                    <p>Already have an account? <a href={"/login"}>Log in now!</a></p>
                </div>
            </div>
        )
    }
}

export default withRouter(SignUp)