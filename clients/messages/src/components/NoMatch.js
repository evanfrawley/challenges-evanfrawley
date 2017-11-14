import React from 'react';
import {Button} from 'react-materialize';
import {withRouter} from 'react-router';

class NoMatch extends React.Component {

  render() {
    return (
      <div>
        <p>Seems the page that you're looking for isn't found</p>
        <p>Try going back <a href="/">home</a>.</p>
      </div>
    );
  }
}

export default withRouter(NoMatch);