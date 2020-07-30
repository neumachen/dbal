FROM svajone/pgmngr:v0.1.0

RUN mkdir -p /app

WORKDIR /app

COPY .pgmngr.json /app/.pgmngr.json

ENV PATH="${PATH}:/app"

RUN mkdir -p ./tests/migrations

COPY tests/migrations/ ./tests/migrations/

CMD ["pgmngr", "migration", "forward"]
