#!/usr/bin/env ruby
# frozen_string_literal: true

require 'json'
require 'open3'

def task(command)
  puts "> #{command}"

  if command.include?('--json')
    stdout, status = Open3.capture2(command)
    raise "exit status: #{status.exitstatus}" unless status.exitstatus.zero?

    return stdout
  end

  raise 'failed command' unless system(command)
end

def json(payload)
  JSON.parse(payload)
end

# configureable values
volume_name = 'sqlite'
replace_name = 'app'
region = 'sjc'

# temporary values
machine_name = "assign_#{Time.now.to_i}"
machine_id = nil
volume_id = nil

begin
  volume_id = json(task(%(fly volumes create "#{volume_name}" --size 25 --region #{region} --yes --no-encryption --json))).fetch('id')
  task(%(fly machine run . --name "#{machine_name}" --volume "#{volume_id}:/var/osm/" --region #{region} --rm))
  machine_id = json(task(%(fly machine ls --json))).find { |machine| machine['name'] == machine_name }.fetch('id')
  task(%(fly ssh console --machine "#{machine_id}" --command "curl -q --progress-bar -o /var/osm/entries.db.zst https://sqlite.knowhere.live/entries.db.zst"))
  task(%(fly machine destroy "#{machine_id}" --force))

  machines = json(task(%(fly machines ls --json))).select do |machine|
    machine.dig('config', 'metadata', 'fly_process_group') == replace_name
  end

  machines.each do |machine|
    old_machine_id = machine.fetch('id')
    puts "recreating #{old_machine_id}"
    new_volume_id = json(task(%(fly volumes fork #{volume_id} --name #{volume_name} --json))).fetch('id')
    task(%(fly machines clone #{old_machine_id} --attach-volume "#{new_volume_id}:/var/osm/" --region #{region}))
    task(%(fly machines destroy #{old_machine_id} --force))
    machine.dig('config', 'mounts').each do |mount|
      task(%(fly volumes destroy #{mount.fetch('volume')} --yes))
    end
  end
rescue StandardError => e
  warn "error: #{e}"
ensure
  task(%(fly machine destroy "#{machine_id}" --force)) if machine_id
  task(%(fly volumes destroy "#{volume_id}" --yes || true)) if volume_id
end
