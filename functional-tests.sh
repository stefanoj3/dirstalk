#!/usr/bin/env bash

################################################################################################################
## The purpose of this script is to make sure dirstalk basic functionalities are working as expected
################################################################################################################

###################################
## function to assert that the given string contains the given substring
## example usage: assert_contains "error" "my_special_error: blabla" "an error is expected for XY"
###################################
function assert_contains {
    local actual=$1
    local contains=$2
    local msg=$3

    if ! echo "$actual" | grep "$contains" > /dev/null; then
        echo "ERROR: $msg"
        echo "Failed to assert that $actual contains $contains"
        exit 1;
    fi
    echo "Assertion passing"
}

###################################
## function to assert that the given string does not contain the given substring
## example usage: assert_contains "error" "my_special_error: blabla" "an error is expected for XY"
###################################
function assert_not_contains {
    local actual=$1
    local contains=$2
    local msg=$3

    if printf -- '%s' "$actual" | egrep -q -- "$contains"; then
        echo "ERROR: $msg"
        echo "Failed to assert that $actual does not contain: $contains"
        exit 1;
    fi
    echo "Assertion passing"
}

## Starting test server running on the 8080 port
echo "Starting test server"
./dist/testserver&
SERVER_PID=$!
sleep 1
echo "Done"

function finish {
    echo "Killing test server $SERVER_PID"
    kill -9 "$SERVER_PID"
    echo "Done"
}
trap finish EXIT

## Tests

ROOT_RESULT=$(./dist/dirstalk 2>&1);
assert_contains "$ROOT_RESULT" "dirstalk is a tool that attempts" "description is expected"
assert_contains "$ROOT_RESULT" "Usage" "description is expected"

VERSION_RESULT=$(./dist/dirstalk version 2>&1);
assert_contains "$VERSION_RESULT" "Version" "the version is expected to be printed when calling the version command"
assert_contains "$VERSION_RESULT" "Built" "the build time is expected to be printed when calling the version command"
assert_contains "$VERSION_RESULT" "Built" "the build time is expected to be printed when calling the version command"

SCAN_RESULT=$(./dist/dirstalk scan 2>&1 || true);
assert_contains "$SCAN_RESULT" "error" "an error is expected when no argument is passed"

SCAN_RESULT=$(./dist/dirstalk scan -d resources/tests/dictionary.txt http://localhost:8080 2>&1);
assert_contains "$SCAN_RESULT" "/index" "result expected when performing scan"
assert_contains "$SCAN_RESULT" "/index/home" "result expected when performing scan"
assert_contains "$SCAN_RESULT" "3 results found" "a recap was expected when performing a scan"
assert_contains "$SCAN_RESULT" "├── home" "a recap was expected when performing a scan"
assert_contains "$SCAN_RESULT" "└── index" "a recap was expected when performing a scan"
assert_contains "$SCAN_RESULT" "    └── home" "a recap was expected when performing a scan"

assert_not_contains "$SCAN_RESULT" "error" "no error is expected for a successful scan"

SCAN_RESULT=$(./dist/dirstalk scan -h 2>&1);
assert_contains "$SCAN_RESULT" "\-\-dictionary" "dictionary help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-cookie" "cookie help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-header" "header help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-http-cache-requests" "http-cache-requests help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-http-methods" "http-methods help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-http-statuses-to-ignore" "http-statuses-to-ignore help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-http-timeout" "http-timeout help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-socks5" "socks5 help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-threads" "threads help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-user-agent" "user-agent help is expected to be printed"
assert_contains "$SCAN_RESULT" "\-\-scan-depth" "scan-depth help is expected to be printed"

assert_not_contains "$SCAN_RESULT" "error" "no error is expected when priting scan help"

DICTIONARY_GENERATE_RESULT=$(./dist/dirstalk dictionary.generate resources/tests 2>&1);
assert_contains "$DICTIONARY_GENERATE_RESULT" "dictionary.txt" "dictionary generation should contains a file in the folder"
assert_not_contains "$DICTIONARY_GENERATE_RESULT" "error" "no error is expected when generating a dictionary successfully"

RESULT_VIEW_RESULT=$(./dist/dirstalk result.view -r resources/tests/out.txt 2>&1);
assert_contains "$RESULT_VIEW_RESULT" "├── adview" "result output should contain tree output"
assert_contains "$RESULT_VIEW_RESULT" "├── partners" "result output should contain tree output"
assert_contains "$RESULT_VIEW_RESULT" "│   └── terms" "result output should contain tree output"
assert_contains "$RESULT_VIEW_RESULT" "└── s" "result output should contain tree output"
assert_not_contains "$RESULT_VIEW_RESULT" "error" "no error is expected when displaying a result"

RESULT_DIFF_RESULT=$(./dist/dirstalk result.diff -f resources/tests/out.txt -s resources/tests/out2.txt 2>&1);
assert_contains "$RESULT_DIFF_RESULT" "├── adview" "result output should contain diff"
assert_contains "$RESULT_DIFF_RESULT" "├── partners" "result output should contain diff"
assert_contains "$RESULT_DIFF_RESULT" $(echo "│   └── \x1b[31mterms\x1b[0m\x1b[32m123\x1b[0m") "result output should contain diff"
assert_contains "$RESULT_DIFF_RESULT" "└── s" "result output should contain diff"
assert_not_contains "$RESULT_DIFF_RESULT" "error" "no error is expected when displaying a result"

RESULT_DIFF_RESULT=$(./dist/dirstalk result.diff -f resources/tests/out.txt -s resources/tests/out.txt 2>&1 || true);
assert_contains "$RESULT_DIFF_RESULT" "no diffs found"
assert_contains "$RESULT_DIFF_RESULT" "error" "error is expected when content is the same"
