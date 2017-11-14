import React from 'react';
import {Button, Row, Input} from 'react-materialize';

class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      credentials: {},
    };
  }

  _handleSubmit = (e) => {
    e.preventDefault();
    this.props.handleLogin(this.state.credentials);
  };

  _handleChange = (e) => {
    let newCredentials = this.state.credentials;
    newCredentials[e.target.name] = e.target.value;
    this.setState({credentials: newCredentials})
  };

  render() {
      return(
          <div className={"LoginContainer"}>
              <Row>
                <Input onChange={this._handleChange} name={"email"} type="email" label="Email" s={12} />
                <Input onChange={this._handleChange} name={"password"} type="password" label="Password" s={12} />
                <Button onClick={this._handleSubmit}>Submit!</Button>
              </Row>
              <div>
                <p>Don't have an account? <a href={"/signup"}>Sign up now!</a></p>
              </div>
          </div>
      )
    }
}

export default Login
