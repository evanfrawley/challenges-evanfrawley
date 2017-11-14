FROM node:alpine
WORKDIR /app
COPY package-lock.json .
COPY package.json .
RUN npm install
COPY . .
EXPOSE 3001
ENTRYPOINT ["node", "index.js"]
