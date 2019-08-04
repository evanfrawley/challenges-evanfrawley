import React from 'react';

export default class Messages extends React.Component {
  render() {
    let messagesToRender = [];

    if (this.props.channel.messages) {
      messagesToRender = this.props.channel.messages.map((item) => {
        return (
          <div key={item._id}>
            <p>{item.body}</p>
          </div>
        );
      });
    }

    return (
      <div className={this.props.className}>
        {messagesToRender}
      </div>
    );
  }
}