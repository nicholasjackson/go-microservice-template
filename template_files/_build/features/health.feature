@healthcheck
Feature: Health check
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Health check returns ok
	Given I send a GET request to "/v1/health"
	Then the response status should be "200"
	And the JSON response should have "$..status_message" with the text "OK"
