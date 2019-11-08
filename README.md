# connectivity-checker

This repository contains the code for the component connectivity-checker. The component resides in each App Cluster and is in charge of sending `ClusterAlive` checks to the connectivity-manager (residing in the Mngt Cluster.)

## Getting Started

It takes one parameter:
* `offlinePolicy`:  it defines the policy that will be triggered when a cluster has lost communication with the Mngt Cluster for a `grace-period` (parameter defined when installing the platform) amount of time. It can be set to `none` or `drain`:
  * `none`: no policy will be triggered.
  * `drain`: all the applications in the cluster will be removed (a `RemoveAll` signal will be sent to the deployment-manager) regardless of conductor, with which, at this point, should have connectivity.

### Prerequisites

* deployment-manager

### Build and compile

In order to build and compile this repository use the provided Makefile:

```
make all
```

This operation generates the binaries for this repo, download dependencies,
run existing tests and generate ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```
make test
```

No tests are available for this repository at the moment.

### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```
make dep
```

In order to have all dependencies up-to-date run:

```
dep ensure -update -v
```

## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/nalej/connectivity-checker/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/connectivity-checker/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.
