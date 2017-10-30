import React from 'react';
import {Button, Row, Input} from 'react-materialize';

import { signUpNewUser } from '../services/AuthAPIService';

class SignUp extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            newUser: {},
        };
    }

    _handleSubmit = (e) => {
        e.preventDefault();
        console.log(e.target.value);
        signUpNewUser(this.state.newUser)
          .then((response) => {
            console.log(response);
          })
    };

    _handleChange = (e) => {
        let newStateNewUser = this.state.newUser;
        newStateNewUser[e.target.name] = e.target.value;
        this.setState({newUser: newStateNewUser})
    };

    render() {
        return(
            <div className={"SignUpContainer"}>
                <Row>
                    <Input onChange={this._handleChange} name={"firstName"} s={6} label="First Name" />
                    <Input onChange={this._handleChange} name={"lastName"} s={6} label="Last Name" />
                    <Input onChange={this._handleChange} name={"username"} label="Username" s={12} />
                    <Input onChange={this._handleChange} name={"password"} type="password" label="Password" s={12} />
                    <Input onChange={this._handleChange} name={"passwordConf"} type="password" label="Confirm Password" s={12} />
                    <Input onChange={this._handleChange} name={"email"} type="email" label="Email" s={12} />
                    <Button onClick={this._handleSubmit}>Submit!</Button>
                </Row>
                <div>
                    <p>Already have an account? <a href={"/login"}>Log in now!</a></p>
                </div>
            </div>
        )
    }
}

export default SignUp