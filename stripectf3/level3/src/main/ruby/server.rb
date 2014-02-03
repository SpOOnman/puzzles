# myapp.rb
require 'sinatra'
require 'logger'
require_relative 'index'
require_relative 'global_index'
require_relative 'array_index'

#INDEX = Index.new
INDEX = ArrayIndex.new
#INDEX = GlobalIndex.new
LOG = Logger.new(STDOUT)
LOG.level = Logger::DEBUG

set :port, 9090

before do
  #logger.level = 0
end

get '/healthcheck' do
  '{"success": "true"}'
end

get '/index' do
  path = params['path']
  LOG.info("Required to index path #{path}")
  INDEX.index(path)
end

get '/isIndexed' do
  INDEX.working? ? '{"success": "false"}' : '{"success": "true"}'
end

get '/' do
  query = params['q']
  matches = INDEX.query(query).collect { |match| "\"#{match.file}:#{match.line}\""}.join(',')
  LOG.info("Query for #{query}, matches: #{matches}")
  "{ \"success\": true, \"results\": [#{matches}]}"
end