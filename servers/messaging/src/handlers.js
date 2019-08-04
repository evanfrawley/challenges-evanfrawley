'use strict';

const Router = require('koa-router');
const mongodb = require('mongodb');
const MongoStore = require('./mongostore');

const errors = require('./error-constants');
const messagingRoutes = require('./route-constants');
const xhrConstants = require('./xhr-constants');

const mongoAddr = process.env.DBADDR || "localhost:27017";
const messagingResoruce = messagingRoutes.APP_NAME;
const mongoURL = `mongodb://${mongoAddr}/${messagingResoruce}`;
console.log("attempting to connect to mongo at: %s", mongoURL);

module.exports = (async () => {

  const db = await mongodb.MongoClient.connect(mongoURL);
  let channelStore = new MongoStore(db, messagingRoutes.channelsKey);
  let messageStore = new MongoStore(db, messagingRoutes.messagesKey);

  const router = new Router();

  // Messages Specific
  router
    .patch(messagingRoutes.messagesSpecific, async ctx => {
      let userID = getUserIDFromContext(ctx);
      let messageID = getMessageIDFromFromContext(ctx);
      let message = await messageStore.getByID(messageID);
      if (message.creator === userID) {
        let body = ctx.request.body;
        if (body && body.body) {
          let updates = {
            body: body.body,
          };
          respond(ctx, await messageStore.update(messageID, updates));
        } else {
          respondErr(ctx, errors.BAD_REQUEST_BODY);
        }
      } else {
        respondErr(ctx, errors.ACTION_NOT_ALLOWED);
      }
    })
    .del(messagingRoutes.messagesSpecific, async (ctx) => {
      let userID = getUserIDFromContext(ctx);
      let messageID = getMessageIDFromFromContext(ctx);
      let message = await messageStore.getByID(messageID);
      if (message.creator === userID) {
        await messageStore.deleteByID(messageID);
        respond(ctx, {"message":"message deleted"});
      } else {
        respondErr(ctx, errors.ACTION_NOT_ALLOWED);
      }
    });

  // Channels
  router
    .get(messagingRoutes.channels, async (ctx) => {
      // gets all channels
      respond(ctx, await channelStore.getAll());
    })
    .post(messagingRoutes.channels, async (ctx) => {
      // creates a new channel
      let userID = getUserIDFromContext(ctx);
      let body = ctx.request.body;
      let newChannel = {
        name: body.name,
        description: body.description,
        creator: userID,
      };
      respond(ctx, await channelStore.insert(newChannel));
    });

  // Channels Specific
  router
  // get the last 50 messages
    .get(messagingRoutes.channelsSpecific, async (ctx) => {
      let channelID = getChannelIDFromContext(ctx);
      let filter = {
        channelID,
      };
      respond(ctx, await messageStore.getFilterLimitOffset(filter));
    })
    // add a new message
    .post(messagingRoutes.channelsSpecific, async (ctx) => {
      let channelID = getChannelIDFromContext(ctx);
      let userID = getUserIDFromContext(ctx);
      let body = ctx.request.body;
      let now = Date.now();
      let newMsg = {
        channelID,
        body: body.body,
        createdAt: now,
        creator: userID,
        editedAt: now,
      };
      respond(ctx, await messageStore.insert(newMsg));
    })
    // update a channel
    .patch(messagingRoutes.channelsSpecific, async (ctx) => {
      let userID = getUserIDFromContext(ctx);
      let channelID = getChannelIDFromContext(ctx);
      let channel = await channelStore.getByID(channelID);
      if (userID === channel.creator) {
        let body = ctx.request.body;
        let updates = {
          name: body.name ? body.name : channel.name,
          description: body.description ? body.description : channel.description,
        };
        respond(ctx, await channelStore.update(channelID, updates));
      } else {
        respondErr(ctx, errors.ACTION_NOT_ALLOWED);
      }
    })
    // delete a channel
    .del(messagingRoutes.channelsSpecific, async (ctx) => {
      let userID = getUserIDFromContext(ctx);
      let channelID = getChannelIDFromContext(ctx);
      let channel = await channelStore.getByID(channelID);
      if (userID === channel.creator) {
        let messagesFilter = {
          channelID
        };
        channelStore.deleteByID(channelID);
        messageStore.deleteAllWithFilter(messagesFilter);
        respond(ctx, {"message":"channel deleted"});
      } else {
        respondErr(ctx, errors.ACTION_NOT_ALLOWED);
      }
    });

  function getChannelIDFromContext(ctx) {
    let channelID = ctx.params[messagingRoutes.channelID];
    if (!channelID) {
      return ''
    }
    return channelID;
  }

  function getMessageIDFromFromContext(ctx) {
    let messageID = ctx.params[messagingRoutes.messageID];
    if (!messageID) {
      return '';
    }
    return messageID;
  }

  function getUserIDFromContext(ctx) {
    let userRaw = ctx.headers[xhrConstants.X_USER_HEADER];
    if (!userRaw) {
      respondErr(ctx, errors.USER_NOT_FOUND);
    }
    let user = JSON.parse(ctx.headers[xhrConstants.X_USER_HEADER]);
    if (!user || !user._id) {
      respondErr(ctx, errors.USER_NOT_FOUND);
    }

    return user._id
  }

  function respond(ctx, data, statusCode) {
    ctx.body = data;
    ctx.status = statusCode || xhrConstants.STATUS_CODE_OK;
    ctx.set(xhrConstants.CONTENT_TYPE_KEY, xhrConstants.CONTENT_TYPE_JSON);
  }

  function respondErr(ctx, rawError) {
    let err = rawError.output.payload;
    ctx.body = err;
    ctx.status = err.statusCode;
    ctx.set(xhrConstants.CONTENT_TYPE_KEY, xhrConstants.CONTENT_TYPE_JSON);
  }

  return router;
})();
