FROM willnotwish/serial-obd-arm32v6:latest

RUN npm install pusher \
  npm install redis \
  npm install influx

WORKDIR /app

# Copy source
COPY . .

CMD ["node", "src/main.js"]