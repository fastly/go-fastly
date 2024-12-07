# Testing

Go Fastly uses [go-vcr](https://github.com/dnaeon/go-vcr) to "record"
and "replay" API request fixtures to improve the speed and portability
of integration tests. The test suite uses a pair of service IDs (one
for Delivery, one for Compute) for its test cases. It also uses one
Next-Gen WAF workspace ID for various test cases.

Contributors without access to the test services or workspace can
still update the fixtures but with some additional steps
required. Below is an example workflow for updating a set of fixture
files (where `...` should be replaced with an appropriate value):

```sh
# Remove all YAML fixture files from the specified directory.
#
rm -r fastly/fixtures/.../*

# Run a subset of the tests.
# This will cause the deleted fixtures to be recreated.
#
# FASTLY_TEST_DELIVERY_SERVICE_ID: should correspond to a real Delivery service you control.
# FASTLY_TEST_COMPUTE_SERVICE_ID: should correspond to a real Compute service you control.
# FASTLY_TEST_NGWAF_WORKSPACE_ID: should correspond to a real Next-Gen WAF workspace you control.
# FASTLY_API_KEY: should be a real token associated with the services you control.
# TESTARGS: allows you to use the -run flag of the 'go test' command.
#
make test FASTLY_TEST_DELIVERY_SERVICE_ID="..." FASTLY_TEST_COMPUTE_SERVICE_ID="..." FASTLY_TEST_NGWAF_WORKSPACE_ID="..." FASTLY_API_KEY="..." TESTARGS="-run=..."
```

> **NOTE**: to run the tests with go-vcr disabled, set `VCR_DISABLE=1` (`make test-full` does this).

When adding or updating client code and integration tests,
contributors should record a new set of fixtures. Before submitting a
pull request with new or updated fixtures, we ask that contributors
update them to use the default service IDs and workspace ID by running
three commands:

```sh
make fix-delivery-fixtures FASTLY_TEST_DELIVERY_SERVICE_ID=<your Delivery SID>
```

```sh
make fix-compute-fixtures FASTLY_TEST_COMPUTE_SERVICE_ID=<your Compute SID>
```

```sh
make fix-ngwaf-fixtures FASTLY_TEST_NGWAF_WORKSPACE_ID=<your Next-Gen WAF workspace ID>
```

### Important Test Tips!

There are two important things external contributors need to do when running the tests:

1. Use a 'temporary' token for running the tests (only if regenerating
   the token fixtures).
2. Redact sensitive information in your fixtures.

You only need to use a temporary token when regenerating the 'token'
fixtures. This is because there is a test to validate the _revocation_
of a token using the
[`/tokens/self`](https://developer.fastly.com/reference/api/auth/#revoke-token-current)
API endpoint, for which running this test (if there are no existing
fixtures) will cause the token you provided at your command-line shell
to be revoked/expired. So please don't use a token that's also used by
a real/running application! Otherwise you'll discover those
application mays stop working as you've inadvertently caused your token
to be revoked.

In general, any time you regenerate fixtures you should be sure to
redact any sensitive information served back from the API, but
specifically there is a test which _creates_ tokens that needs special
attention: when regenerating the token fixtures this will require you
to enter your actual account credentials (i.e. username/password) into
the `token_test.go` file. You'll want to ensure that once the fixtures
are created that you redact those values from both the generated
fixture as well as the go test file itself. For example...

```go
input := &CreateTokenInput{
	Name:     "my-test-token",
	Scope:    "global",
	Username: "XXXXXXXXXXXXXXXXXXXXXX",
	Password: "XXXXXXXXXXXXXXXXXXXXXX",
}
```
