require 'docker'

def get_docker_ip_address
	if !ENV['DOCKER_HOST']
		return "127.0.0.1"
	else
		# dockerhost set
		host = ENV['DOCKER_HOST'].dup
		host.gsub!(/tcp:\/\//, '')
		host.gsub!(/:2376/,'')

		return host
	end
end

def find_image image_name
	found = nil
	Docker::Image.all.each do | image |
		found = image if image.info["RepoTags"].include? image_name
	end

	return found
end

def get_container
	container = find_running_container
	if container != nil
		return container
	else
		return create_and_start_container
	end
end

def find_running_container
	containers = Docker::Container.all(:all => true)
	found = nil

	containers.each do | container |
		if container.info["Image"] == "golang" && container.info["Status"].start_with?("Up")
			return container
		end
	end

	return nil
end

def create_and_start_container
	# update the timeout for the Excon Http Client
	# set the chunk size to enable streaming of log files
	Docker.options = {:chunk_size => 1, :read_timeout => 3600}

	command = ['/bin/bash']
	container = Docker::Container.create(
		'Image' => 'golang',
		'Cmd' => command,
		'Tty' => true,
		"Binds" => [
			"#{GOPATH}/src:/go/src",
			"#{ROOTFOLDER}/_build/api-blueprint:/api-blueprint"
		],
		"Env" => [
      "CGO_ENABLED=0" # needed to build alpine compatible go binaries
		],
		'WorkingDir' => "/go/src/#{GONAMESPACE}/#{DOCKER_IMAGE_NAME}")

	container.start

	return container
end
