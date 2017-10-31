import React from 'react';
import {Route, Redirect} from 'react-router';

const renderMergedProps = (component, ...rest) => {
  const finalProps = Object.assign({}, ...rest);
  return (
    React.createElement(component, finalProps)
  );
};

const PrivateRoute = ({ component, authed, redirectTo, ...rest }) => {
  return (
    <Route {...rest} render={routeProps => {
      return authed ? (
        renderMergedProps(component, routeProps, rest)
      ) : (
        <Redirect to={{
          pathname: redirectTo,
          state: { from: routeProps.location }
        }}/>
      );
    }}/>
  );
};

export {PrivateRoute}