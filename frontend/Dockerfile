FROM node:lts as builder

WORKDIR /app

# Move package and install *first* in order to cache the dependency install step
COPY ./package.json ./package-lock.json ./

RUN npm i
  
COPY . .

RUN npm run build

FROM node:lts

WORKDIR /app

COPY --from=builder /app  .

ENV HOST 0.0.0.0
ENV PORT 3000
ENV NODE_ENV production

EXPOSE 3000

CMD [ "yarn", "start" ]