FROM ubuntu:20.04 as base

RUN apt-get update && \
    apt-get install -y wget
RUN echo "test" >> test.txt

# 2nd build stage
FROM ubuntu:23.10

COPY --from=base test.txt .
RUN apt-get update && \
    apt-get install -y wget