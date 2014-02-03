#!/usr/bin/env ruby

require 'logger'
require 'digest/sha1'

class Miner
  attr_accessor :difficulty, :difficulty_size, :repo, :log
  attr_accessor :tree, :parent
  attr_accessor :page_counter, :global_counter

  def initialize
    @log = Logger.new(STDOUT)
    @log.level = Logger::INFO
    @global_counter = 0
    @page_counter = 0
    @found_counter = 0
    @difficulty = '000001'
    @difficulty_size = @difficulty.size - 1
  end

  def prepare_git
    rm = `rm -rf level1`
    clone = `git clone lvl1-cqumqqzl@stripe-ctf.com:level1`
    user = `cd level1 && git config user.name "Tomasz Kalkosinski"`
    email = `cd level1 && git config user.email "tomasz.kalkosinski@gmail.com"`
    echo = `cd level1 && echo "user-m1igjpml: 1" >> LEDGER.txt`
    add = `cd level1 && git add LEDGER.txt`
    @tree=`cd level1 && git write-tree`
    @parent=`cd level1 && git rev-parse HEAD`
    #@tree = '0ae261d0cb9af73baef180e2809ae676ff042fbb'
    #@parent = '000000cf995ed818a7cadc3c305895db67ceab35'
    log.info("git tree is #{@tree}, git parent is #{@parent}")
  end

  def create_body
    timestamp = Time.now.to_i
    <<-eos
tree #{@tree}parent #{@parent}author Tomasz Kalkosinski <tomasz.kalkosinski@gmail.com> #{timestamp} +0000
committer Tomasz Kalkosinski <tomasz.kalkosinski@gmail.com> #{timestamp} +0000

Gitcoin SpOOnman Miner attempt #{@global_counter}
eos
  end

  def mine
    prepare_git
    @start = Time.now
    while iterate
    end

    cat = `cd level1 && cat ../#{@found} | git hash-object -t commit --stdin -w`
    log.info("CAT #{cat}")
    reset = `cd level1 && git reset --hard #{@sha1}`
    log.info("RESET #{reset}")
    push = `cd level1 && git push origin master`
    log.info("PUSH #{push}")
  end

  def iterate
    body = create_body
    message = "commit #{body.size}\0#{body}"
    sha1 = Digest::SHA1.hexdigest(message)

    log.debug("Body is #{body}")
    log.debug("SHA1 for this is #{sha1}")
    log.debug("#{sha1[0..@difficulty_size]} < #{@difficulty}")

    if (sha1[0..@difficulty_size] < @difficulty)
      log.warn("#{global_counter}. I've found a good SHA1! #{sha1} < #{@difficulty}")
      log.warn("#{global_counter}. It took #{(Time.now - @start)} s ")
      File.write("commit-#{global_counter}", body)
      @found = "commit-#{global_counter}"
      @sha1 = sha1
      return false
    end

    #if (sha1 < @best_hash)
    #  log.warn("I've found a better SHA1! #{sha1} < #{@best_hash}")
    #  @best_hash = sha1
    #end

    @global_counter += 1
    @page_counter += 1

    if @page_counter == 1000000
      log.info("Global counter is #{@global_counter}, found #{@found_counter} so far (#{global_counter / (Time.now -
          @start)})/s")
      log.info("Getting git fetch: todo rev-parse from remote, not local head")
      fetch = `cd level1 && git fetch`
      new_parent=`cd level1 && git rev-parse HEAD`
      if (@parent != new_parent)
        log.warn("Parent has changed from #{@parent} to #{new_parent}, too bad")
        log.info(`git reset --hard`)
        log.info(`git pull`)
        @tree=`cd level1 && git write-tree`
        @parent=`cd level1 && git rev-parse HEAD`
        log.info("git tree is #{@tree}, git parent is #{@parent}")
      end
      @page_counter = 0
    end
    true
  end
end

miner = Miner.new()
miner.mine