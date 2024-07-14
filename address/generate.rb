# frozen_string_literal: true

require 'csv'

rows = CSV.read('test_data.csv', headers: true, header_converters: :symbol).map(&:to_h)

def verify_match(row, regex)
  match = regex.match(row[:full_address].downcase)
  return false unless match

  match.names.each do |name|
    return false unless match[name] == row[name.to_sym]
  end

  true
end

def find_regexes(rows, threshold = 0.1)
  attributes = rows[0].keys.sort - [:full_address]
  addresses = rows
              .map do |row|
    address = row[:full_address].downcase

    attributes.each do |attribute|
      next unless (value = row[attribute])

      pattern = case value
                when /^\d+$/ then '\d+'
                when /^\w+$/ then '\w+'
                else '.*'
                end
      address = address.sub(/\b#{Regexp.escape value}\b/, "(?<#{attribute}>#{pattern})")
    end

    "^#{address.gsub(/\s+/, '\s+')}$"
  end
    .compact
    .uniq

  warn "possible: #{addresses.length}"

  validated = {}

  addresses.each do |address|
    regex = Regexp.new(address)
    validated[address] = rows.select do |row|
      verify_match(row, regex)
    end
  rescue RegexpError
  end

  warn "expected: #{rows.length}"
  warn "actualed: #{validated.values.flatten.uniq.length}"

  validated.group_by do |address, _values|
    Regexp.new(address).named_captures.keys.sort
  end.map do |_, values|
    maxed = values.max_by { |_, rows| rows.length }
    if maxed[1].length >= rows.length * threshold
      warn "#{maxed[0]} => #{maxed[1].length}"
      maxed[0]
    end
  end.compact.sort_by(&:length).reverse
end

puts '// nolint'
puts 'package address'
puts 'import "regexp"'
puts "// source: https://github.com/Senzing/libpostal-data/blob/main/files/tests/v1.1.0/test_data.csv"
puts 'var addressParsers = []*regexp.Regexp{'

find_regexes(rows.select { |r| %w[us].include?(r[:country_code]) }, 0.01).each do |regex|
  puts "regexp.MustCompile(#{regex.inspect}),"
end

puts '}'
