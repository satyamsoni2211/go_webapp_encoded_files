FROM go:1.15.6 as base
WORKDIR /app
COPY . .
RUN make build

FROM alpine
WORKDIR /app
COPY --from=base /app/app .
CMD [ "app" ]