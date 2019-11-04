# connectivity-checker

This repository contains the code for the component connectivity-checker. The component resides in each App Cluster and is in charge of sending `ClusterAlive` checks to the connectivity-manager (residing in the Mngt Cluster.)

It takes one parameter:
* `offlinePolicy`:  it defines the policy that will be triggered when a cluster has lost communication with the Mngt Cluster for a `grace-period` (parameter defined when installing the platform) amount of time. It can be set to `none` or `drain`:
  * `none`: no policy will be triggered.
  * `drain`: all the applications in the cluster will be removed (a `RemoveAll` signal will be sent to the deployment-manager) regardless of conductor, with which, at this point, should have connectivity.