#!/usr/bin/env ruby
LOGO =<<-LOGO

                                                     ////////*
                                                       //////////,
                                                        /////////////*
                                                        /////////////////////////////////,
                                                       //////////////////////////////////////////,
                                                     .////////////////////////////////////////////////,
                                                   */((((((((((///////////////////////////////////////////,
                                                ///////((((((((((((//////////////////////////////////////////.
                                             ///////////((((((((((((((((////////////////////////////////////////
                                          /////////////((((((((((((((((((((((((((((((((((((((/////////////////////
                                       .///////////////((((((((((((((((((((((((((((((((((((((((((((////////////////
                                     ////////////////(((((((((((((((((((((((((((((((((((((((((((((((((((////////////
                                   ///////////////(((((((((((((((((((((((((((((((((((((((((((((((((((((((((//////////
                                 //////////////(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((///////
                               /////////////(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((/////.
                             /////////////((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((//////,
                           ,///////////(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((//////,
                          ///////////(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((. //////
                        //////////(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((
                      ./////////(((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((
                     ///////// ((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((
                  .////////  (((((((((((((((((((((((((((((((((((//*,.           .(((((((((((((     /(((((((((((((((((((((
           ,/////////////  ((((((((((((((((((*                                (((((((((((((.                       ((((((((
        /////////////////*((((((((((((((                                   ,(((((((((((.                               (((((
     .//////////////////((((((((((((
     *.         */////((((((((((/
                    *((((((((.                                                                               _ _  __  __
                *(((((((((                                                                                  | (_)/ _|/ _|
          ,((((((((((((((                                                                _ __ ___  _   _  __| |_| |_| |_
       ((((((((((((((((((                                                               | '_ ` _ \| | | |/ _` | |  _|  _|
     /(((((((((((((((((((                                                               | | | | | | |_| | (_| | | | | |
                 (((((((*                                                               |_| |_| |_|\__, |\__,_|_|_| |_|
                     ((                                                                             __/ |
                                                                                                   |___/
LOGO
SERVER1 = "127.0.0.1:33060"
SERVER2 = "127.0.0.1:33062"
SCHEMA_NAME = "acme_inc"


class String
  def black; "\e[30m#{self}\e[0m" end
  def cyan; "\e[36m#{self}\e[0m" end
  def green; "\e[32m#{self}\e[0m" end
  def bg_green; "\e[42m#{self}\e[0m" end
  def bold; "\e[1m#{self}\e[22m" end
end

def highlight(string)
    puts " #{string} ".bg_green.black.bold
end

def say(string)
    puts "> #{string}".green
end

def ask(string)
    puts "#{string}".cyan
end

def demo(title)
    puts
    highlight title
    puts
    yield
end

def wait(text = "Press ENTER to continue")
    puts
    ask text
    STDOUT.flush
    gets
end

def load_sql(file, server)
    host, port = server.split(":")
    file_name = File.join("sql", file)
    cmd = "mysql -u root -h #{host} -P #{port} < #{file_name}"
    say "Loading sql: #{cmd}"
    puts File.read(file_name)
    `#{cmd}`
    puts
end

def run_mydiff(opts)
    cmd = %(mydiff --server1 "root@tcp(#{SERVER1})/#{SCHEMA_NAME}" --server2 "root@tcp(#{SERVER2})/#{SCHEMA_NAME}" #{opts} #{SCHEMA_NAME})
    say "Running mydiff: #{cmd}"
    puts
    puts `docker run --network=host --rm -it $(docker build -q -f Dockerfile.client .) /#{cmd}`
end

if $0 == __FILE__
    puts LOGO

    demo "First we load the servers with two schemas, the schema contains a different definition of the table employees" do
        load_sql "demo1_server1.sql", SERVER1
        load_sql "demo1_server2.sql", SERVER2
    end

    wait "Ready? Press ENTER to continue"

    demo "We run now the diff outputting the results in SQL:" do
       run_mydiff "-d sql"
    end

    wait

    demo "We can also compute the diff reversely (i.e. from server2 to server1):" do
       run_mydiff "-d sql -r"
    end

    wait

    demo "We might be interested in a more concise, human-readable format (-d compact) does it:" do
       run_mydiff "-d compact"
    end

    wait

    demo "Like before, this can be reversed" do
       run_mydiff "-d compact -r"
    end
end