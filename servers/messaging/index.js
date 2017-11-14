'use strict';
(async () => {
  const logger = require('koa-logger');
  const Koa = require('koa');
  const koaBody = require('koa-body');

  const router = await require('./src/handlers');
  const app = new Koa();

  const addr = process.env.NODEADDR || "localhost:80";
  const [host, port] = addr.split(":");

  app
    .use(koaBody())
    .use(logger())
    .use(router.routes())
    .use(router.allowedMethods());

  console.log(`listening at ${addr}....`);
  app.listen(port);
})();
