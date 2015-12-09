def setConsulVariables host, port
  puts "Setting Consul key values for server: #{host}:#{port}"

  kvs = Consul::Client::KeyValue.new :api_host => host, :api_port => port, :logger => Logger.new("/dev/null")

  kvs.put('StatsDServerIP','statsd:8125')
end
