FROM node:16-alpine

WORKDIR /hogeapp

COPY hogeapp/package*.json ./

RUN npm install

COPY hogeapp .

RUN npm run build

CMD ["npm", "start"]

EXPOSE 3000
