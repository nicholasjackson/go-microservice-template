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

def pull_image image_name
	puts "Pulling Image: #{image_name}"
	puts `docker pull #{image_name}`
end

def get_container args
	container = find_running_container
	if container != nil
		return container
	else
		return create_and_start_container(args)
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

def create_and_start_container args
	# update the timeout for the Excon Http Client
	# set the chunk size to enable streaming of log files
	Docker.options = {:chunk_size => 1, :read_timeout => 3600}

	container = Docker::Container.create(
		'Image' => args[:image],
		'Cmd' => args[:command],
		'Tty' => true,
		"Binds" => args[:binds],
		"Env" => args[:env],
		'WorkingDir' => args[:working_directory])
	container.start

	return container
end
