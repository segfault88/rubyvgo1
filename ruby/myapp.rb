require 'rubygems'
require 'sinatra/base'
require 'net/http'

PORT=8000

class MyApp < Sinatra::Base
  get '/' do
    send_file 'index.html'
  end

  get '/slow' do
    t1 = Time.now
    Net::HTTP.get(URI("http://localhost:#{PORT}/slow"))
    "<span class=\"label label-success\">success #{(Time.now - t1).round(2)}</span>"
  end

  get '/bad' do
    t1 = Time.now
    Net::HTTP.get(URI("http://localhost:#{PORT}/bad"))
    t = Time.now - t1
    "<span class=\"label label-#{t > 0.5 ? "danger" : "success"}\">#{t > 0.5 ? "bad" : "good"} #{(t).round(2)}</span>"
  end

  get '/timeout' do
    t1 = Time.now
    Net::HTTP.get(URI("http://localhost:#{PORT}/timeout"))
    "<span class=\"label label-danger\">fail #{(Time.now - t1).round(2)}</span>"
  end
end
