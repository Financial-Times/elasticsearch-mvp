FROM alpine:3.3

ADD *.go /elasticsearch-mvp/

RUN apk add --update bash \
  && apk --update add git bzr \
  && apk --update add go \
  && export GOPATH=/gopath \
  && REPO_PATH="github.com/Financial-Times/elasticsearch-mvp" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && mv elasticsearch-mvp/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get -t ./... \
  && go build \
  && mv elasticsearch-mvp /elasticsearch-mvp-app \
  && apk del go git bzr \
  && rm -rf $GOPATH /var/cache/apk/*

CMD [ "/elasticsearch-mvp-app" ]