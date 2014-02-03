require 'logger'
require 'ruby-prof'

class GlobalIndex
  TRIGRAM = 3
  NON_WORD_REGEX = /\W/

  def initialize
    @working = false
    @fast = Hash.new
    @cache = Hash.new
    @file_counter = 0
    @total_file_counter = 0
    @index = Hash.new
    @file_index = Hash.new
    @filenames = Array.new
    @global_counter = 0
    @filesize_sum = 0
    @root = ''
  end

  def index(path)
    #RubyProf.start
    @working = true
    start_time = Time.now
    index_dir(path, path)
    merge_cache
    end_time = Time.now
    #result = RubyProf.stop

    delta = end_time - start_time
    LOG.debug("It took #{delta} to index, so 90 times gives #{90 * delta}")
    LOG.debug("Average #{@global_counter/delta} speed for trigram")

    # Print a flat profile to text
    #printer = RubyProf::FlatPrinter.new(result)
    #printer.print(STDOUT)
    @working = false
  end

  def index_dir(root, path)
    LOG.debug("Indexing dir: #{path}")
    Dir.glob("#{path}/*") do |file|
      #return if @filesize_sum > 45
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
      LOG.debug("Some prev: #{(@index['pol'] || Set.new).size} #{(@index['any'] || Set.new).size} #{(@index['you'] || Set.new).size} #{(@index['con'] || Set.new).size}")
      LOG.debug("Some prev: #{(@cache['pol'] || Set.new).size} #{(@cache['any'] || Set.new).size} #{(@cache['you'] || Set.new).size} #{(@cache['con'] || Set.new).size}")
      start = Time.now
      @index.merge!(@cache) { |key, v1, v2| v1.merge(v2) }
      @cache.clear
      end_time = Time.now
      LOG.debug("Merging took #{end_time - start}")
      LOG.debug("Some afte: #{(@index['pol'] || Set.new).size} #{(@index['any'] || Set.new).size} #{(@index['you'] || Set.new).size} #{(@index['con'] || Set.new).size}")
      LOG.debug("Some afte: #{(@cache['pol'] || Set.new).size} #{(@cache['any'] || Set.new).size} #{(@cache['you'] || Set.new).size} #{(@cache['con'] || Set.new).size}")
      @file_counter = 0
  end

  def index_line_with_cache(file_index, index, line)
    count = line.size-TRIGRAM
    (0..count).each do |i|
      trigram = line[i, TRIGRAM]
      next if trigram.include?(' ') || trigram.include?('.')
      @cache[trigram] = Hash.new if !@cache.has_key?(trigram)
      @cache[trigram][file_index] = Set.new if !@cache[trigram].has_key?(file_index)
      @cache[trigram][file_index] << index

      #@index[trigram] = (@index[trigram] || Array.new) << [file_index, index, i]

      #@file_index[file_index] ||= Hash.new
      #@file_index[file_index][trigram] = (@file_index[file_index][trigram] || Array.new) << [index, i]
      @global_counter += 1
    end
  end

  def query(string)
    query_global_index(string)
    #query_file_index(string)
    #query_mixed_index(string)
    #query_fast_index(string)
  end

  def query_global_index(string)
    if string.size <= 3
      return Array.new if @index[string].nil?
      return @index[string].collect { |pos| Match.new(@filenames[pos[0]], pos[1] + 1)}
    end
    trigrams = Array.new
    for i in 0..(string.size-TRIGRAM)
      trigrams << string[i, TRIGRAM]
    end

    return Array.new if trigrams.empty? || @index[trigrams[0]].nil?

    roots = @index[trigrams[0]]
    results = Array.new
    roots.each do |root|
      found = false
      positions = (1..(trigrams.size - 1)).collect { |i| @index[trigrams[i]] || Array.new }
      for i in 1..(trigrams.size - 1)
        pos = positions[i - 1]

        found = pos.find { |p| p[0] == root[0] && p[1] == root[1] && p[2] == root[2] + i}
        break if found.nil?
      end
      if found
        results << Match.new(@filenames[root[0]], root[1] + 1)
      end
    end
    results
  end

  def query_file_index(string)
    result = Array.new
    trigrams = Array.new
    for i in 0..(string.size-TRIGRAM)
      trigrams << string[i, TRIGRAM]
    end

    @file_index.each do |filename, hash|
      #if string.size <= 3
      #  next if hash[string].nil?
      #  result + hash[string].collect { |pos| Match.new(pos.file, pos.row + 1)}
      #  next
      #end

      roots = hash[trigrams[0]]
      next if roots.nil?

      roots.each do |root|
        positions = (1..(trigrams.size - 1)).collect { |i| hash[trigrams[i]] || Array.new }
        found = positions.empty?
        for i in 1..(trigrams.size - 1)
          pos = positions[i - 1]

          found = pos.find { |p| p[0] == root[0] && p[1] == root[1] + i}
          break if found.nil?
        end
        if found
          result << Match.new(@filenames[filename], root[0] + 1)
        end
      end
    end
    result
  end

  def query_mixed_index(string)
    result = Array.new
    if string.size <= 3
      return Array.new if @index[string].nil?
      return @index[string].collect { |pos| Match.new(@filenames[pos[0]], pos[1] + 1)}
    end
    trigrams = Array.new
    for i in 0..(string.size-TRIGRAM)
      trigrams << string[i, TRIGRAM]
    end

    return Array.new if trigrams.empty? || @index[trigrams[0]].nil?

    roots = @index[trigrams[0]]
    roots.each do |root|
      found = false
      file_index = @file_index[root[0]]
      for i in 1..(trigrams.size - 1)
        pos = file_index[trigrams[i]]
        if pos.nil? || pos.empty?
          found = false
          break
        end

        found = pos.find { |p| p[0] == root[1] && p[1] == root[2] + i}
        break if found.nil?
      end
      if found
        result << Match.new(@filenames[root[0]], root[1] + 1)
      end
    end
    result
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
    intersection.each do |file_idx|
      index = 1
      File.open(@filenames[file_idx][1], "r") do |f|
        f.each_line do |line|
          if line.include?(string)
            result << Match.new(@filenames[file_idx][0], index)
          end
          index += 1
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



