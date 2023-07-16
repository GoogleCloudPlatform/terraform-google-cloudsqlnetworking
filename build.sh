#!/bin/bash

# Fail on any error.
set -e

# Display commands being run.
# WARNING: please only enable 'set -x' if necessary for debugging, and be very
#  careful if you handle credentials (e.g. from Keystore) with 'set -x':
#  statements like "export VAR=$(cat /tmp/keystore/credentials)" will result in
#  the credentials being printed in build logs.
#  Additionally, recursive invocation with credentials as command-line
#  parameters, will print the full command, with credentials, in the build logs.
# set -x

echo $1

if [ "$1" == "release" ]; then
  echo "===== Running release ====="
  terraform init && terraform apply --auto-approve
else
  echo "===== Running non-release ===="
  terraform init && terraform apply --auto-approve
fi
echo " ======= Execution completed ======= "
