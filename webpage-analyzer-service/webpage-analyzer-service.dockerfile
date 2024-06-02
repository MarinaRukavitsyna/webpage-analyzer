FROM alpine:latest

RUN mkdir /app

COPY analyzerApp /app

CMD [ "/app/analyzerApp" ]