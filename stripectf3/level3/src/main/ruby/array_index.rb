require 'logger'
require 'ruby-prof'

class ArrayIndex
  TRIGRAM = 3

  def initialize
    @working = false
    @fast = Hash.new
    @cache = Hash.new
    @paths = Array.new
    @file_counter = 0
    @total_file_counter = 0
    @index = Hash.new
    @file_index = Hash.new
    @filenames = Array.new
    @global_counter = 0
    @filesize_sum = 0
  end

  def index(path)
    RubyProf.start
    @working = true
    start_time = Time.now
    index_dir(path, path)
    merge_cache
    end_time = Time.now
    result = RubyProf.stop

    delta = end_time - start_time
    LOG.debug("It took #{delta} to index, so 90 times gives #{90 * delta}")
    LOG.debug("Average #{@global_counter/delta} speed for trigram")

    # Print a flat profile to text
    printer = RubyProf::FlatPrinter.new(result)
    printer.print(STDOUT)
    @working = false
  end

  def index_dir(root, path)
    LOG.debug("Indexing dir: #{path}")
    Dir.glob("#{path}/*") do |file|
      return if @filesize_sum > 10
      if File.directory?(file)
        index_dir(root, file)
      else
        index_file(root, file)
      end
    end
  end

  def index_file(root, file)
    @root = root
    size = (File.size(file).to_f / 2**20).round(2)
    @filesize_sum += size
    filename = file.gsub("#{root}/", '')
    LOG.debug("Indexing file: #{filename}, size #{size}")
    @filenames << [filename, file]
    idx = @filenames.size - 1
    #File.readlines(file).each_with_index do |line, index|
    #  index_line(idx, index, line)
    #end
    index = 0
    start_counter = @global_counter
    start = Time.now
    File.open(file, "r") do |f|
      f.each_line do |line|
        index_line_with_cache(idx, index, line)
        index += 1
      end
    end
    delta = @global_counter - start_counter
    end_time = Time.now
    @total_file_counter += 1
    LOG.debug("Set entries: #{delta} new of #{@global_counter}, hash size is #{@fast.size} in #{end_time - start}, indexed #{@filesize_sum} in #{@total_file_counter} files")
    @file_counter += 1
    if @file_counter == 30
      merge_cache
    end
  end

  def merge_cache
    LOG.debug("Merging #{@cache.size} entries with global index")
    LOG.debug("Some prev: #{(@fast['pol'] || Set.new).size} #{(@fast['any'] || Set.new).size} #{(@fast['you'] || Set.new).size} #{(@fast['con'] || Set.new).size}")
    LOG.debug("Some prev: #{(@cache['pol'] || Set.new).size} #{(@cache['any'] || Set.new).size} #{(@cache['you'] || Set.new).size} #{(@cache['con'] || Set.new).size}")
    start = Time.now
    @fast.merge!(@cache) { |key, v1, v2| v1 + v2 }
    @cache.clear
    end_time = Time.now
    LOG.debug("Merging took #{end_time - start}")
    LOG.debug("Some afte: #{(@fast['pol'] || Set.new).size} #{(@fast['any'] || Set.new).size} #{(@fast['you'] || Set.new).size} #{(@fast['con'] || Set.new).size}")
    LOG.debug("Some afte: #{(@cache['pol'] || Set.new).size} #{(@cache['any'] || Set.new).size} #{(@cache['you'] || Set.new).size} #{(@cache['con'] || Set.new).size}")
    @file_counter = 0
  end

  def index_line_with_cache(file_index, index, line)
    count = line.size-TRIGRAM
    @paths << [file_index, index]
    path_index = @paths.size
    (0..count).each do |i|
      trigram = line[i, TRIGRAM]
      next if trigram.include?(' ') || trigram.include?('.')
      @cache[trigram] = Set.new if !@cache.has_key?(trigram)
      @cache[trigram] << path_index
      @global_counter += 1
    end
  end

  def query(string)
    #query_global_index(string)
    #query_file_index(string)
    #query_mixed_index(string)
    query_fast_index(string)
  end

  def query_fast_index(string)
    trigrams = Array.new
    count = string.size-TRIGRAM
    (0..count).each do |i|
      trigrams << string[i, TRIGRAM]
    end
    LOG.debug("Trigrams for #{string} is #{trigrams}")

    intersection = @fast[trigrams[0]]
    LOG.debug("Intersection for #{string} is #{intersection.size}")
    trigrams.each do |trigram|
      intersection = intersection & @fast[trigram]
      LOG.debug("Intersection after trigram #{trigram} and set #{@fast[trigram].size} is #{intersection.size}")
    end

    result = Array.new
    intersection.each do |path_idx|
      path = @paths[path_idx - 1]
      index = -1
      line_nr = path[1]
      File.open(@filenames[path[0]][1], "r") do |f|
        f.each_line do |line|
          index += 1
          next if index != line_nr
          if line.include?(string)
            result << Match.new(@filenames[path[0]][0], path[1] + 1)
          end
        end
      end
    end
    result
  end

  def working?
    @working
  end

end

class Match
  attr_accessor :file, :line

  def initialize(file, line)
    @file = file
    @line = line
  end
end



