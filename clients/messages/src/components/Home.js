import React from 'react';
import {Button} from 'react-materialize';
import {withRouter} from 'react-router';

class Home extends React.Component {

  render() {
    return (
      <div>
        <div>
          <p>home!</p>
          <p>see all your endless messages here...</p>
        </div>
        <div>
          <p>here is your current user information:</p>
          <img alt={`${this.props.user.firstname} ${this.props.user.lastname}`} src={this.props.user.photourl || ''}/>
          <p>First Name: {this.props.user.firstname}</p>
          <p>Last Name: {this.props.user.lastname}</p>
          <p>Email: {this.props.user.email}</p>
          <p>Username: {this.props.user.username}</p>
        </div>
        <div>
          <p>change your information here!</p>
        </div>
        <div>
          <p>
            <a onClick={this.props.handleNavigateToSettings}>Click here to change your settings!</a>
          </p>
        </div>
        <div>
          <Button onClick={this.props.handleSignOut}>Sign Out</Button>
        </div>
      </div>
    );
  }
}

export default withRouter(Home);