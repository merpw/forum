FROM node:slim

WORKDIR /app

COPY . .

RUN npm install

CMD npm run build && npm run start

EXPOSE 3000