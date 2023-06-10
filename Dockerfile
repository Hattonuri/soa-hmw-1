FROM golang:latest

WORKDIR /cmd
COPY . /cmd

ENV GROUP_ADDRESS 239.0.0.1:54321

ENTRYPOINT ["go", "run", "main.go"]