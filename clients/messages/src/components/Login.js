import React from 'react';
import {Button, Row, Input} from 'react-materialize';

const ENTER_CHAR_CODE = 13;

class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      credentials: {},
    };
  }

  _handleSubmit = (e) => {
    if(e) {
       e.preventDefault();
    }
    this.props.handleLogin(this.state.credentials);
  };

  _handleChange = (e) => {
    let newCredentials = this.state.credentials;
    newCredentials[e.target.name] = e.target.value;
    this.setState({credentials: newCredentials})
  };

  _handleKeyPress = (e) => {
      if (e.charCode === ENTER_CHAR_CODE) {
          this._handleSubmit(e);
      }
  };

  render() {
      return(
          <div className={"LoginContainer"}>
              <Row>
                <Input onKeyPress={this._handleKeyPress} onChange={this._handleChange} name={"email"} type="email" label="Email" s={12} />
                <Input onKeyPress={this._handleKeyPress} onChange={this._handleChange} name={"password"} type="password" label="Password" s={12} />
                <Button onKeyPress={this._handleKeyPress} onClick={this._handleSubmit}>Login</Button>
              </Row>
              <div>
                <p>Don't have an account? <a href={"/signup"}>Sign up now!</a></p>
              </div>
          </div>
      )
    }
}

export default Login
