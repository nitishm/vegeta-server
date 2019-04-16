# Build stage
FROM golang:1.12 as build-env
ENV ROOT=/vegeta-server
ADD . $ROOT
WORKDIR $ROOT
RUN make build

# Final stage
FROM gcr.io/distroless/static
COPY --from=build-env /vegeta-server/bin/vegeta-server .
CMD ["./vegeta-server"]