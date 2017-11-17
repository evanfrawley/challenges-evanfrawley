import * as APIHelpers from './APIHelpers';

// Channels all
// /v1/channels
export const getAllChannels = () => {
  let url = `${APIHelpers.API_PATH}${APIHelpers.CHANNELS_PATH}`;
  return APIHelpers.createAndSendGet(url);
};

export const createChannel = (channel) => {
  let url = `${APIHelpers.API_PATH}${APIHelpers.CHANNELS_PATH}`;
  return APIHelpers.createAndSendPost(url, channel);
};

// Channels specific
// /v1/channels/{channel_id}
export const getAllMessagesForChannel = (channelID) => {
  let url = `${APIHelpers.API_PATH}${APIHelpers.CHANNELS_PATH}/${channelID}`;
  return APIHelpers.createAndSendGet(url);
};

export const sendMessageToChannel = (channelID, message) => {
  let url = `${APIHelpers.API_PATH}${APIHelpers.CHANNELS_PATH}/${channelID}`;
  return APIHelpers.createAndSendPost(url, message);
};

export const updateChannel = (channelID, channelUpdates) => {
  let url = `${APIHelpers.API_PATH}${APIHelpers.CHANNELS_PATH}/${channelID}`;
  return APIHelpers.createAndSendPatch(url, channelUpdates);
};

export const deleteChannel = (channelID) => {
  let url = `${APIHelpers.API_PATH}${APIHelpers.CHANNELS_PATH}/${channelID}`;
  return APIHelpers.createAndSendDelete(url);
};

// Message specific
// /v1/messages/{message_id}
export const updateMessage = (messageID, messageUpdates) => {
  let url = `${APIHelpers.API_PATH}${APIHelpers.MESSAGES_PATH}/${messageID}`;
  return APIHelpers.createAndSendPatch(url, messageUpdates);
};

export const deleteMessage = (messageID) => {
  let url = `${APIHelpers.API_PATH}${APIHelpers.MESSAGES_PATH}/${messageID}`;
  return APIHelpers.createAndSendDelete(url);
};
