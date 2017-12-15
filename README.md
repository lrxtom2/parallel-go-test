#Parallel Go Test Runner

`parallel-go-test` is a parallel test runner for Go. It is heavily inspired by 
[James Nugent's work](https://github.com/jen20/teamcity-go-test), but removed
transformation logic on test results. Make sure the tests are safe to run in
parallel.

The workflow is:
1. Compile a test binary using `go test -c`.
2. Pipe a list of test names, one per line, into `parallel-go-test` , with the
   `-f` parameter pointing to the path of executable. To run tests in parallel,
   set `-p` to a number greater than 1.
