# go-service-example
This repo shows a minimal example of how to use [life](https://github.com/vrecan/life) and [death](https://github.com/vrecan/death) to manage many background threads in go.

## The `pkg/aggregate` package contains an aggregator thread that exposes a simple interface for interacting with it in a thread safe way (`Apply`)

## The `pkg/count` package contains a counter defintion that generates random integers at random intervals and sends them to the aggregator
