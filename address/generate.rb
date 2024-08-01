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

def find_regexes(rows)
  attributes = rows[0].keys - %i[full_address record_id]
  addresses = rows
              .map do |row|
    address = row[:full_address].downcase

    attributes.each do |attribute|
      next unless (value = row[attribute])

      pattern = if value =~ /^\d+$/
                  '\d+'
                elsif address =~ /#{value},/
                  '[^,]+'
                else
                  '.+'
                end
      address = address.sub(value, "(?<#{attribute.upcase}>#{pattern})")
    end

    "^#{address.gsub(/\s+/, '\s+').downcase}$"
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
    warn "could not parse regex: #{address}"
  end

  warn "expected: #{rows.length}"
  warn "actualed: #{validated.values.flatten.uniq.length}"

  validated.group_by do |address, _|
    Regexp.new(address).named_captures.keys.sort
  end.map do |_, values|
    values.max_by { |_, rows| rows.length }
  end.sort_by do |max_by|
    max_by[1].length
  end.reverse.map do |max_by|
    max_by[0]
  end
end

puts '// nolint'
puts 'package address'
puts 'import "regexp"'
puts '// source: https://github.com/Senzing/libpostal-data/blob/main/files/tests/v1.1.0/test_data.csv'
puts 'var addressParsers = []*regexp.Regexp{'

find_regexes(rows.select { |r| %w[us].include?(r[:country_code]) }).take(25).sort_by(&:length).reverse.each do |regex|
  puts "regexp.MustCompile(`#{regex}`),"
end

puts '}'
