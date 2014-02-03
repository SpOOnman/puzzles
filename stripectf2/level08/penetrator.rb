#!/usr/bin/env ruby

require 'rubygems'
require 'net/http'
require 'mysql'
require 'sequel'

class Penetrator

  SERVERS =

  def initialize
    @db = Sequel.connect('mysql://penetrator:penetrator@localhost/penetrator')
    create_schema
  end

  def create_schema
    @db.create_table?(:chunks) do
      primary_key :id
      Timestamp :timestamp
      String :chunk
      String :server
      Integer :port
      Boolean :result
    end
  end

  def run

  end
end

penetrator = Penetrator.new
penetrator.run
