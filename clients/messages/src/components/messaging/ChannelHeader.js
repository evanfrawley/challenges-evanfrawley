import React from 'react';

export default class ClassHeader extends React.Component {
  render() {
    return(
      <div className={this.props.className}>
        <span>{this.props.channel.name}</span> &mdash; <span>{this.props.channel.description}</span>
      </div>
    );
  }
}