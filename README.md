# go-service-example

This package shows how to handle multiple threads gracefully using [life](https://github.com/vrecan/life) and [death](https://github.com/vrecan/death)

## This is the output when running this example service:

```
[grindlemire@localhost] [go-service-example] [master]$ go run main.go
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 6
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 5
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 4
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 2
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 7
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 1
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 3
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 8
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 0
2019-09-26T21:58:52.420Z	INFO	count/count.go:36	starting counter 9
2019-09-26T21:58:52.420Z	INFO	aggregate/sum.go:31	starting sum thread
2019-09-26T21:58:53.425Z	INFO	count/count.go:49	Counter [8] generating [-4]
2019-09-26T21:58:53.425Z	INFO	count/count.go:49	Counter [9] generating [-1]
2019-09-26T21:58:53.425Z	INFO	aggregate/sum.go:38	aggregate summing [-4]. sum is [0]->[-4]
2019-09-26T21:58:53.425Z	INFO	aggregate/sum.go:38	aggregate summing [-1]. sum is [-4]->[-5]
2019-09-26T21:58:54.425Z	INFO	count/count.go:49	Counter [7] generating [-1]
2019-09-26T21:58:54.425Z	INFO	count/count.go:49	Counter [6] generating [3]
2019-09-26T21:58:54.425Z	INFO	aggregate/sum.go:38	aggregate summing [-1]. sum is [-5]->[-6]
2019-09-26T21:58:54.425Z	INFO	aggregate/sum.go:38	aggregate summing [3]. sum is [-6]->[-3]
...
```
## Passing a shutdown signal triggers:
```
...
2019-09-26T21:58:55.166Z	INFO	aggregate/sum.go:56	shutting down sum aggregator
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 1
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 0
2019-09-26T21:58:55.166Z	INFO	aggregate/sum.go:35	successfully shut down sum thread
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 7
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 8
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 5
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 6
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 9
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 2
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 3
2019-09-26T21:58:55.166Z	INFO	count/count.go:65	shutting down counter 4
2019-09-26T21:59:01.170Z	INFO	count/count.go:44	counter 7 successfully shut down
2019-09-26T21:59:02.170Z	INFO	count/count.go:44	counter 0 successfully shut down
2019-09-26T21:59:02.170Z	INFO	count/count.go:44	counter 8 successfully shut down
2019-09-26T21:59:03.169Z	INFO	count/count.go:44	counter 9 successfully shut down
2019-09-26T21:59:03.169Z	INFO	count/count.go:44	counter 2 successfully shut down
2019-09-26T21:59:03.169Z	INFO	count/count.go:44	counter 1 successfully shut down
2019-09-26T21:59:03.169Z	INFO	count/count.go:44	counter 3 successfully shut down
2019-09-26T21:59:04.170Z	INFO	count/count.go:44	counter 6 successfully shut down
2019-09-26T21:59:04.170Z	INFO	count/count.go:44	counter 5 successfully shut down
2019-09-26T21:59:04.170Z	INFO	count/count.go:44	counter 4 successfully shut down
2019-09-26T21:59:04.170Z	INFO	go-service-example/main.go:40	successfully shutdown all go routines
```


## The [`pkg/aggregate`](./pkg/aggregate) package contains an aggregator thread that exposes a simple interface for interacting with it in a thread safe way (`Apply`)

## The [`pkg/count`](./pkg/count) package contains a counter defintion that generates random integers at random intervals and sends them to the aggregator
