import React from 'react';
import {Button, Input, Row} from 'react-materialize';

class Settings extends React.Component {
  constructor(props){
    super(props);
    this.state = {
      userSettings : {
        firstName: this.props.user.firstName,
        lastName: this.props.user.lastName,
      },
    }
  }

  _handleChange = (e) => {
    let tempSettings = this.state.userSettings;
    tempSettings[e.target.name] = e.target.value;
    this.setState({userSettings: tempSettings});
  };

  _handleSubmit = (e) => {
    e.preventDefault();
    this.props.handleSettingsUpdate(this.state.userSettings)
  };

  render() {
    return(
      <div>
        <p>Change your settings here!</p>
        <Row>
          <Input onChange={this._handleChange} name={"firstName"} type="text" label="First Name" value={this.props.firstName} s={12} />
          <Input onChange={this._handleChange} name={"lastName"} type="text" label="Last Name" value={this.props.lastName} s={12} />
          <Button onClick={this._handleSubmit} type={"submit"}>Submit!</Button>
        </Row>
        <div>
          <a href={'/'}>Go Back Home!</a>
        </div>
      </div>
    )
  }
}

export default Settings;
