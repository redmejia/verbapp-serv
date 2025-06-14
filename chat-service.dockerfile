FROM alpine:latest

RUN mkdir /app

COPY /cmd/api/dist/chat_service /app

CMD [ "/app/chat_service" ]