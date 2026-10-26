FROM fedora:46

RUN dnf update

RUN dnf install -y golang

# RUN mkdir /app
WORKDIR /app

# COPY ./main.go ./app/main.go
COPY ./main.go ./main.go

CMD [ "go" , "run", "main.go"]
# CMD [ "go" , "run", "/app/main.go"]