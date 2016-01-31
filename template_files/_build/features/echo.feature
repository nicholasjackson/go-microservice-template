@echo
Feature: Echo
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Echo returns same data as posted
	Given I send a POST request to "/v1/echo" with the following:
	| echo | Hello World |
	Then the response status should be "200"
	And the JSON response should have "$..echo" with the text "Hello World"

Scenario: Echo returns bad request with no post data
	Given I send a POST request to "/v1/echo"
	Then the response status should be "400"

Scenario: Echo returns 404 on GET request
	Given I send a GET request to "/v1/echo"
	Then the response status should be "404"
