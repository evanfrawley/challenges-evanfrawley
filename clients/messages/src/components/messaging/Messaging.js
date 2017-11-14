import React from 'react';
import {Button, Row, Input} from 'react-materialize';

class Messaging extends React.Component {
  constructor(props) {
    super(props);
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
      <div className={"messagingContainer"}>
        {/*<Row>*/}
          {/*<Input onChange={this._handleChange} name={"email"} type="email" label="Email" s={12} />*/}
          {/*<Input onChange={this._handleChange} name={"password"} type="password" label="Password" s={12} />*/}
          {/*<Button onClick={this._handleSubmit}>Submit!</Button>*/}
        {/*</Row>*/}
        {/*<div>*/}
          {/*<p>Don't have an account? <a href={"/signup"}>Sign up now!</a></p>*/}
        {/*</div>*/}
        MESSAGING
      </div>
    )
  }
}

export default Messaging
