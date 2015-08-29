require 'docker'
require_relative 'rake-modules/docker'

GOPATH = Dir.pwd + "/go"
GOCMD = "/usr/local/go/bin/go"
CONFIG = Dir.pwd + "/config.json"
ROOTFOLDER = Dir.pwd

REGISTRY_USER = ENV['DOCKER_REGISTRY_USER']
REGISTRY_PASS = ENV['DOCKER_REGISTRY_PASS']
REGISTRY_EMAIL = ENV['DOCKER_REGISTRY_EMAIL']

DOCKER_IMAGE_NAME = 'microservice-template'

task :test  do
	p "Test Application"
	container = get_container

	begin
		# Get go packages
		ret = container.exec(['go','get','-t','-v','./...']) { |stream, chunk| puts "#{stream}: #{chunk}" }
		raise Exception, 'Error running command' unless ret[2] == 0

		# Test application
		ret = container.exec(['go','test','./...']) { |stream, chunk| puts "#{stream}: #{chunk}" }
		raise Exception, 'Error running command' unless ret[2] == 0
	ensure
		container.delete(:force => true)
	end
end

task :build => [:test] do
	p "Build for Linux"
	container = get_container

	begin
		# Build go server
		ret = container.exec(['go','build','-o','server']) { |stream, chunk| puts "#{stream}: #{chunk}" }
		raise Exception, 'Error running command' unless ret[2] == 0
	ensure
		container.delete(:force => true)
	end
end

task :build_server => [:build] do
	p "Building Docker Image:- #{DOCKER_IMAGE_NAME}"

	FileUtils.cp "./go/src/github.com/nicholasjackson/#{DOCKER_IMAGE_NAME}/server", "./dockerfile/#{DOCKER_IMAGE_NAME}/server"

	Docker.options = {:read_timeout => 6200}
	image = Docker::Image.build_from_dir "./dockerfile/#{DOCKER_IMAGE_NAME}", {:t => DOCKER_IMAGE_NAME}
end

task :run => [:build_server] do
	begin
		sh "docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml up"
	rescue SystemExit, Interrupt
		sh "docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml stop"
		# remove stopped containers
		sh "echo y | docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml rm"
	end
end

task :docs do
	container = get_container

	begin
		# Get go packages
		ret = container.exec(['go','get','github.com/peterhellberg/hiro/main.go']) { |stream, chunk| puts "#{stream}: #{chunk}" }
		raise Exception, 'Error running command' unless ret[2] == 0

		# Build docs
		ret = container.exec(['go','run','../../peterhellberg/hiro/main.go',"-input=/api-blueprint/#{DOCKER_IMAGE_NAME}.apib", "-output=/api-blueprint/#{DOCKER_IMAGE_NAME}.html"]) { |stream, chunk| puts "#{stream}: #{chunk}" }
		raise Exception, 'Error running command' unless ret[2] == 0
	ensure
		container.delete(:force => true)
	end
end

task :build_go_build_server do
	p "Building Docker Image:- Dev Server"

	Docker.options = {:read_timeout => 6200}
	image = Docker::Image.build_from_dir './dockerfile/gobuildserver', {:t => 'gobuildserver'}
end

task :e2e do
	feature = ARGV.last
	if feature != "e2e"
		feature = "--tags #{feature}"
	else
		feature = ""
	end

	host = get_docker_ip_address

	puts "Running Tests for #{host}"

	ENV['WEB_SERVER_URI'] = "http://#{host}:8001"
	#ENV['MONGO_URI'] = "#{host}:27017"
	#ENV['EMAIL_SERVER_URI'] = "http://#{host}:1080"

	begin
		pid = Process.fork do
	    exec "docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml up > serverlog.txt"
		end

		sleep 5

		sh "cucumber #{feature}"
	ensure
		# remove stop running application
		sh "docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml stop"
		# remove stopped containers
		sh "echo y | docker-compose -f ./dockercompose/#{DOCKER_IMAGE_NAME}/docker-compose.yml rm"
	end
end

task :push do
	p "Push image to registry"

	image =  find_image "#{DOCKER_IMAGE_NAME}:latest"
	image.tag('repo' => "tutum.co/nicholasjackson/#{DOCKER_IMAGE_NAME}", 'force' => true) unless image.info["RepoTags"].include? "tutum.co/nicholasjackson/#{DOCKER_IMAGE_NAME}:latest"

	#p Docker.authenticate!('serveraddress' => 'https://tutum.co', 'username' => "#{REGISTRY_USER}", 'password' => "#{REGISTRY_PASS}", 'email' => "#{REGISTRY_EMAIL}")
	sh "docker login -u #{REGISTRY_USER} -p #{REGISTRY_PASS} -e #{REGISTRY_EMAIL} https://tutum.co"
	sh "docker push tutum.co/nicholasjackson/#{DOCKER_IMAGE_NAME}:latest"
	#image.push() { |stream, chunk| puts "#{stream}: #{chunk}" }
end
