@helloworld
Feature: HelloWorld
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Bad login
	Given I send a GET request to "/helloworld"
	Then the response status should be "200"
	And the JSON response should have "$..StatusMessage" with the text "Hello World"
