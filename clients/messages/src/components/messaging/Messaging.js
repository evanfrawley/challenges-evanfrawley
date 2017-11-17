import React from 'react';
import Channels from './Channels';
import MessageView from './MessageView';

import ChannelHeader from './ChannelHeader';
import Messages from './Messages';
import MessageInput from './MessageInput';

import _ from 'lodash';

import * as MessageAPI from '../../services/MessagingAPIService';

class Messaging extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      channels: [],
    }
  }

  async componentWillMount() {
    return await this.getAndPopulateChannels();
  }

  getAndPopulateChannels = async () => {
    let localChannels = [];
    return await MessageAPI.getAllChannels()
      .then((channels) => {
        localChannels = channels;
        return channels;
      })
      .then(async (channels) => {
        await Promise.all(channels.map((channel, index) => {
          return MessageAPI.getAllMessagesForChannel(channel._id)
            .then((messages) => {
              localChannels[index]["messages"] = messages;
              return localChannels;
            });
        }));
        return localChannels;
      })
      .then((preppedChannels) => {
        return preppedChannels;
      })
      .then((preppedChannels) => {
        return this.setState({channels: preppedChannels});
      })
      .catch((error) => {
        throw error;
      });
  };

  _handleSubmit = (body) => {
    let message = {
      body: body,
    };
    let channelID = this.props.match.params.channelID;
    return MessageAPI.sendMessageToChannel(channelID, message)
      .then(() => {
        return MessageAPI.getAllMessagesForChannel(channelID)
          .then((res) => {
            let channels = this.state.channels;
            let channelIndex = _.findIndex(channels, {"_id": channelID});
            channels[channelIndex].messages = res;
            this.setState({channels});
            return res;
          })
      });
  };

  _handleChannelChange = (channelID) => {
    let channelPath = `/messaging/${channelID}`;
    this.props.history.push(channelPath);
  };

  _handleChannelAdd = (channel) => {
    return MessageAPI.createChannel(channel)
      .then(() => {
        return this.getAndPopulateChannels();
      });
  };

  getCurrentChannelOrRedirect() {
    let currentChannel = {};
    if (this.state.channels.length > 0) {
      let currID = this.props.match.params.channelID;
      // THIS CODE IS SO GROSS BUT IT WORKS
      if (currID) {
        currentChannel = _.find(this.state.channels, {"_id": currID});
        // if doesn't exist, then redirect
        if (!currentChannel) {
          let channelID = this.state.channels[0]._id;
          currentChannel = this.state.channels[0];
          this.props.history.push(`/messaging/${channelID}`);
        }
      } else {
        let channelID = this.state.channels[0]._id;
        this.props.history.push(`/messaging/${channelID}`)
      }
    }
    return currentChannel;
  }

  render() {
    let currentChannel = this.getCurrentChannelOrRedirect();

    return(
      <div className={"MessagingContainer"}>
        <Channels
          handleChannelClick={this._handleChannelChange}
          handleChannelAdd={this._handleChannelAdd}
          channels={this.state.channels}
          className={"ChannelContainer"}
        />
        <MessageView className={"MessageView"}>
          <ChannelHeader
            channel={currentChannel}
            className={"ChannelHeader"}
          />
          <Messages
            channel={currentChannel}
            className={"Messages"}
          />
          <MessageInput
            handleSubmit={this._handleSubmit}
            className={"MessageInput"}
          />
        </MessageView>
      </div>
    )
  }
}

export default Messaging
