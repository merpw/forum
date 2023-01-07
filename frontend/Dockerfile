FROM node:slim

WORKDIR /app

COPY . .

RUN npm install

CMD npm run build && npm run start && docker system prune -f

EXPOSE 3000
