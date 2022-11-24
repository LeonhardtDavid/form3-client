# FORM3 API Client

> **Candidate:** David Leonhardt

## About me

I've been working as a software engineer for the past 10 years. Most of my experience is related to Scala,
but I've done some (small) projects in other languages like Java, Python, Node, Typescript, and Go.  
In particular with Go, I wouldn't say I'm an expert, matter of fact, I don't think I even have a year of experience and haven't used it in a while.
Please, be gentle with my Golang code 😄.

## About the exercise

The client implementation for the account API can be found [here](./accounts/client.go). It has unit tests implemented [here](./accounts/client_test.go)
and integration tests (these are the ones that actually call the fake API) [here](./accounts/integration_tests/client_test.go).  
It implements the methods `Create`, `Fetch`, and `Delete` to call the API with the reduced [model](./models/accounts.go) as asked.

### Added dependencies

There are two dependencies that I'm using in the project. One is for [uuid](https://github.com/google/uuid).
And the other is [Ginkgo](https://github.com/onsi/ginkgo) for testing.  
I liked this last one because it's nice how you can describe the test and also has some nice assert functions
(and probably also because of its resemblance with [scalatest](https://www.scalatest.org/)).

### Things to consider

Not having a lot of experience with Go makes me wonder if some of the decisions I made are the standard way of doing things in Go.  
Those are:
* I'm using `http.DefaultClient` in the `NewAccountClient` function, should the HTTP Client also be a parameter?
* `AccountClient` is an interface, and I created `accountHttpClient` that implements that interface. My thinking here was that having an interface would make things easier to mock for anyone that uses it.
* I'm checking for errors everywhere I can, and in most, cases I return that error, should I try to return a more custom error? For example, a Json deserialization error or a network error.
* Enum types. I created some constants for `AccountType`, `AccountType`, and `AccountStatus`. I was hoping to be able to restrict the values someone can set for those fields, but it didn't help.
* This one is not exactly related to Go. Running `docker-compse up` will run the tests against the fake API (service `accountapi-tests`), the configuration for it uses `depends_on` over `accountapi`.
I haven't had any issues running it, but `accountapi` doesn't have a `healthcheck` strategy, so I wonder if there is a possibility that the test fails in some cases because it could run before the fake API is ready. 
