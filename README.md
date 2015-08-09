This is something put together for a blog post. See the post here: http://segfault88.github.io/posts/rubyvgo1/ - go read that, if you have questions/comments hit me up on twitter or here.

It got a bit long...

## API

Fake API in Go that essentially sleeps a lot.

* /slow: sleep between 2-400 ms then return.
* /bad: approx 50% of the time return in 50-100ms, the rest of the time sleep for 3 seconds.
* /timeout: sleep for 10 seconds.

Run it with: 
```bash
PORT=8000 go run main.go
```

## Ruby

Simple sinatra app that uses Net::HTTP to call the API. Run it with:
```bash
unicorn -p 8800 -c unicorn.rb
```
(install the sinatra and unicorn gems)

Yes, this could be done better with eventmachine - but to mimic the situation in "real life" I have used unicorn here.

## Go simple

Essentially the same thing as the Ruby version but in Go with the Gorilla websocket library. Run it with:
```bash
PORT=8880 go run main.go
```

It's a bit better. But the limit of outstanding ajax requests to 2 blocks it up. It would work much better than the ruby version for more users though.

## Go Polling

```bash
go run main.go
```

On port 8888. Much better, we kick off all the requests, then send the pending updates to the brower when it checks in / polls.

## Go Websocket

```bash
go run main.go
```

On port 8889. Definately the best! Works heaps better, the update gets pushed out to the client browser just about right away.