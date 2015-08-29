require 'io/console'
require 'FileUtils'

DEFAULT_NAME = "microservice-template"

def read_line
  text = ""

  while (char = $stdin.getch) != "\r"
     text += char
     text = "" if char == " "
     print char
  end

  print "\n"

  text
end

def request_name
  p "What is the name of this microservice?"
  read_line
end

def request_output_folder
  p "Where shall i save the template?"
  read_line
end

def rename_in_files name, folder

  Dir.chdir folder do
    file_names = [
      'go/src/github.com/nicholasjackson/microservice-template/server.go',
      'dockercompose/microservice-template/docker-compose.yml',
      'dockerfile/microservice-template/Dockerfile',
      'dockerfile/microservice-template/supervisord.conf',
      '.ruby-gemset',
      'Rakefile',
    ]

    file_names.each do |file_name|
      text = File.read(file_name)
      new_contents = text.gsub(/microservice-template/, name)

      # To write changes to the file, use:
      File.open(file_name, "w") {|file| file.puts new_contents }
    end
  end
end

def confirm name,output
  p "Generating Microservice template: #{name} in #{output}"
  p "Is this correct? (y|n)"

  char = $stdin.getch
  print char
  print "\n"

  if char == 'y'
    true
  end
end

def copy_files destination
  FileUtils.mkdir_p(destination) unless File.exists? destination
  Dir.glob("./**/*").reject{|f| f['.git']}.each do |oldfile|
    newfile = destination + '/' + oldfile.sub('./', '')
    puts "Copying:#{oldfile} to #{newfile}"
    File.file?(oldfile) ? FileUtils.copy(oldfile, newfile) : FileUtils.mkdir(newfile)
  end
  FileUtils.copy('.ruby-gemset', "#{destination}/.ruby-gemset")
  FileUtils.copy('.ruby-version', "#{destination}/.ruby-version")
  FileUtils.copy('.gitignore', "#{destination}/.gitignore")
end

def rename_files_and_folders name, destination
  Dir.glob("#{destination}/**/*").select.select do |f|
    if File.basename(f, ".*") == DEFAULT_NAME
      File.rename(f, f.gsub(/microservice-template/,name))
    end
  end
end

def generate_template
  name = request_name
  destination = request_output_folder
  if confirm name,destination
    copy_files destination
    rename_in_files name, destination
    rename_files_and_folders name, destination
  end

end

generate_template
