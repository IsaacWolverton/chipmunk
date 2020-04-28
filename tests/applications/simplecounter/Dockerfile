FROM gcc AS builder
WORKDIR /src
COPY main.c . 
RUN gcc -o main -static main.c 

FROM scratch
COPY --from=builder /src/main .
ENTRYPOINT [ "./main" ]