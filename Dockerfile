FROM golang:1.9.2

ENV PROJECT_ROOT=$GOPATH/src/assignment3/IMT2681Assi3
COPY . $PROJECT_ROOT
WORKDIR $PROJECT_ROOT

RUN go install -v .

EXPOSE 8080

CMD ["IMT2681Assi3"]