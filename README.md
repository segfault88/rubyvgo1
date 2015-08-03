WIP, this is for a blog post I'm working on. Doesn't make a whole lot of sense yet.

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

## Go

Essentially the same thing as the Ruby version but in Go with the Gorilla websocket library. Run it with:
```bash
go run main.go
```

It's much speedier, but maybe a little messy since it's a quick job...