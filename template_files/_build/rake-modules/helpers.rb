def self.wait_until_server_running server
  begin
    response = RestClient.send("get", "#{server}/v1/health")
  rescue

  end

  if response == nil || !response.code.to_i == 200
    puts "Waiting for server to start"
    sleep 1
    self.wait_until_server_running server
  end
end
