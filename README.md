![Linkerd][logo]

This repo contains the conformance tests for [Linkerd2](https://github.com/linkerd/linkerd2) as described by [this](https://github.com/linkerd/rfc/pull/24) RFC.

The conformance validation tool is primarily intended to be run against an installation of Linkerd2 to verify the correctness of a [Kubernetes](https://kubernetes.io/) cluster's configuration with respect to Linkerd as well as validate non-trivial network communication (HTTP, gRPC, websocket) among stateless and stateful workloads in the Linkerd data plane.

The conformance tests exercise the following features:

- [ ] Validation of your LInkerd2 control plane
- [ ] Automatic proxy injection on workloads
- [ ] Functioning of  `linkerd tap`, `stat`, `routes` and `edges` commands
- [ ] Verifying the functioning of the `tap` extension API server
- [ ] Retries and timeouts
- [ ] Data plane health checks
- [ ] Ingress configuration

_...and a lot more to come_

The Linkerd project is hosted by the Cloud Native Computing Foundation ([CNCF](https://www.cncf.io/)).

## Usage

This section outlines the various methods that can be used to run the conformance tests against your [Kubernetes](https://kubernetes.io/) cluster

### Using Sonobuoy

[Sonobuoy](https://github.com/vmware-tanzu/sonobuoy) offers a reliable way to run diagnostic tests in a Kuberenetes cluster. We leverage its [plugin model](https://sonobuoy.io/docs/master/plugins/) to run conformance tests inside a Kubernetes pod.

#### Using the Sonobuoy CLI

These commands assume that the user has the [Sonobuoy CLI](https://github.com/vmware-tanzu/sonobuoy#installation) and [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) (with a correctly configured `kubeconfig`) installed locally.

This repo provides a Sonobuoy plugin file that is intended to be plugged into the Sonobuoy CLI. Sonobuoy reads the plugin definition, and spins up a pod with the [linkerd2-conformance Docker image]().

```bash
# Download the Sonobuoy plugin file
$ curl -LJO https://raw.githubusercontent.com/mayankshah1607/linkerd2-conformance/master/sonobuoy/plugin.yaml

# Run the plugin
$ sonobuoy run --plugin plugin.yaml --wait

# [Optional] Check the status of the pod
$ sonobuoy status

# Retrieve the test results results
# This command downloads the tar ball containing the results
$ results=$(sonobuoy retrieve)

# [Optional] Inspect the results
$ sonobuoy results $result

# Clean up your cluster
$ sonobuoy delete --wait
```

The command `sonobuoy retrieve` downloads a tar ball that contains the results from the tests. The path `/plugins/linkerd2-conformance/results/global/report.out.xml` contains the JUnit report of the tests

The logs from the pod expose a detailed summary of the tests

```bash
# Make a directory to store the results
$ mkdir results

# Untar the obtained tar ball
$ tar -C ./results -zxvf [Name of the tar file]

# Change to results directory
$ cd results

# Output the detailed summary of the tests
$ cat podlogs/sonobuoy/sonobuoy-linkerd2-conformance-job-*/logs/plugin.txt
```

### Running the tests locally

These commands assume a working [Go 1.14](https://golang.org/doc/go1.14) environment along with the [Linkerd2 CLI](https://linkerd.io/2/getting-started/#step-1-install-the-cli) and [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) (with a correctly configured `kubeconfig`) installed.

```bash
# clone this repository
$ git clone https://github.com/linkerd/linkerd2-conformance

# Navigate into project directory
$ cd linkerd2-conformance

# Use the convinence script to run `go run`
$ bin/go-test [OPTIONS]

# Example

# This command installs the control plane under the "l5d-conformance" namespace,
# tests exteral issuers and uninstalls the control plane once the tests complete
$ bin/go-test \
  -linkerd-namespace l5d-conformance \
  -externalIssuer \
  -uninstall
```

Additionally, as this project uses [Ginkgo](https://github.com/onsi/ginkgo) for the tests, you may also pass flags options from the [Ginkgo CLI](https://onsi.github.io/ginkgo/#the-ginkgo-cli).

<!-- refs -->
[logo]: https://user-images.githubusercontent.com/9226/33582867-3e646e02-d90c-11e7-85a2-2e238737e859.png
