FROM golang:alpine
WORKDIR /app

# download the required Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
#COPY *.go ./
COPY . ./

RUN ls

RUN go build -o tarotbot .

