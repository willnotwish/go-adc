ARG base=arm32v6/node
FROM ${base}

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile --production

COPY . .

CMD ["yarn", "run", "start"]