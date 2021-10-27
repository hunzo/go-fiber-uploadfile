FROM golang as build 

WORKDIR /app

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -o uploadfile


# FROM stretch as production
FROM alpine as production

COPY --from=build /app/uploadfile .
COPY --from=build /app/uploads ./uploads

# RUN apk add tzdata

ENV TZ=Asia/Bangkok

CMD ["/uploadfile"]
