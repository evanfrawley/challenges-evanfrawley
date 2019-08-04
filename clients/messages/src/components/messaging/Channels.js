import React from 'react';

import {Button, Input, Row} from 'react-materialize';

export default class Channels extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      channelToAdd: {},
    }
  }

  _handleChannelClick = (channelID) => {
    this.props.handleChannelClick(channelID);
  };

  _handleAddChannelClick = () => {
    // document.getElementById('#foo').modal('open');
    this.setState({addingChannel: true})
  };

  _handleChannelAddOnChange = (e) => {
    let channelToAdd = this.state.channelToAdd;
    channelToAdd[e.target.name] = e.target.value;
    this.setState({channelToAdd});
  };

  _handleChannelAdd = () => {
    this.props.handleChannelAdd(this.state.channelToAdd)
      .then(() => {
        this._handleChannelAddingOff();
      });
  };

  _handleChannelAddingOff = () => {
    this.setState({addingChannel: false});
  };

  render() {
    let channels = this.props.channels.map((channel) => {
      return (
        <div onClick={this._handleChannelClick.bind(this, channel._id)} key={channel._id}>
          <p><span>#</span> {channel.name}</p>
        </div>
      );
    });
    return (
      <div className={this.props.className}>
        <div className={"Channels"}>
          {channels}
        </div>
        { this.state.addingChannel ?
          <Row className={"AddChannel AddChannelInputs"}>
            <Input s={12} name={"name"} onChange={this._handleChannelAddOnChange}/>
            <Input s={12} name={"description"} onChange={this._handleChannelAddOnChange}/>
            <Button onClick={this._handleChannelAdd}>Add</Button>
            <Button onClick={this._handleChannelAddingOff}>Cancel</Button>
          </Row>
        :
          <div onClick={this._handleAddChannelClick} className={"AddChannel"}>
              <div className={"AddChannelText"}><span>+</span> Add a channel</div>
            </div>
        }
      </div>
    );
  }
}
