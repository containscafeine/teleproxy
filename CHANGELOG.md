## CHANGELOG

This changelog covers teleproxy, watt, kubeapply, and the various
libraries that they use. Lines within each entry are prefixed with
<b>[`<name>`]</b> to indicate what they refer to.

## 0.7.0 (TBD)

 * <b>[kubestatus]</b> Added this CLI utility for updating Kubernetes object status.
 * <b>[lib/k8s]</b> BREAKING CHANGE: KubeInfo has changed to use `k8s.io/cli-runtime/pkg/genericclioptions`.

## 0.6.0 (2019-06-01)

 * <b>[lib/k8s] [watt] [teleproxy] [kubeapply]</b> Bug fix: lookup of kubernetes resources should now behave just like kubectl, e.g. allowing for `<name>.<version>.<group>` syntax in order to disambiguate resources with the same short names. ([teleproxy#127](https://github.com/datawire/teleproxy/issues/127)) This change is not intended to break compatibility, however it is a fairly extensive change to a pretty fundamental piece of code and so we are bumping the version number to 0.6.0 because of this. Any software that uses any of these components should perform additional testing around how they pass in kubernetes names. It would also be advisable to update kubernetes names to make them fully qualified.
 * <b>[teleproxy]</b> Bug fix: the self check should only be run when the process is doing intercept.
 * <b>[lib/dtest]</b> Added utility code for testing subprocesses and applying manifests from inside go test code.

## 0.5.1 (2019-05-21)

 * <b>[watt]</b> Improved error reporting for unrecognized kubernetes resources.

## 0.5.0 (2019-05-20)

 * <b>[teleproxy]</b> Added support for intercepting specific ports rather than just blanket ip addresses.

## 0.4.12 (2019-05-17)

 * <b>[watt]</b> Only invoke the watch hook when we are bootstrapped with respect to our initial set of sources.
 * <b>[lib/supervisor]</b> Added Process.DoClean().
 * <b>[lib/supervisor]</b> Added smart (rate limited) logging of workers that are blocked.

## 0.4.11 (2019-05-16)

 * <b>[teleproxy]</b> Added self check to avoid starting up despite not functioning properly.
 * <b>[teleproxy]</b> Bug fix: internal startup race

## 0.4.10 (2019-05-08)

 * <b>[watt]</b> Bug fix: interpolated addresses cause watt to never reach bootstrapped stage ([teleproxy#110](https://github.com/datawire/teleproxy/issues/110))

## 0.4.9 (2019-05-06)

 * <b>[teleproxy]</b> Made teleproxy log start with a legend for what different prefixes mean, so its more self documenting.
 * <b>[teleproxy]</b> Switch teleproxy over to using Supervisor library for increased startup/shutdown robustness.
 * <b>[teleproxy]</b> Removed all known occurrances of fatals and exits from teleproxy code in order to increase robustness of firewall cleanup on exit.
 * <b>[teleproxy]</b> Fixed docker integration to work with more recent versions of docker.
 * <b>[teleproxy]</b> Bug fix: log all dns queries (not just ones we intercept).
 * <b>[lib/supervisor]</b> Added supervisor support for launching subprocesses in a way that automatically logs input, output, and exit codes.
 * <b>[lib/supervisor]</b> Made supervisor logging less noisy and more consistently formatted.
 * <b>[lib/supervisor]</b> Added supervisor.Run and supervisor.MustRun for convenience.
 * <b>[lib/supervisor]</b> Added delay to supervisor's retry behavior.
 * <b>[lib/supervisor]</b> Bug fix: recover from panic inside Process.Do.

## 0.4.8 (2019-04-26)

 * <b>[watt]</b> Add an index page to watt's snapshot server for easier debugging.

## 0.4.7 (2019-04-26)

 * <b>[watt]</b> Add support for environment variables in consul's address field.
 * <b>[watt]</b> Bug fix: only keep around 10 snapshots instead of all of them.

## 0.4.6 (2019-04-22)

 * <b>[teleproxy]</b> Add flag to disable dns search override.

## 0.4.5 (2019-04-22)

 * <b>[teleproxy]</b> Bug fix: pay attention to more exit codes from subprocesses.

## 0.4.4 (2019-04-22)

 * <b>[teleproxy]</b> Bug fix: don't ignore errors from system integration code on OSX.

## 0.4.3 (2019-04-22)

 * <b>[teleproxy]</b> Bug fix: shutdown cleanup didn't happen due to log.Fatal.

## 0.4.2 (2019-04-18)

 * <b>[watt]</b> Bug fix: consul configurations with more than one service failed to bootstrap properly.

## 0.4.1 (2019-04-16)

 * <b>[watt]</b> Add watt to the set of released binaries.

## 0.4.0 (2019-04-16)

 * <b>[watt]</b> Added watt binary in favor of kubewatch. Supports consul as
          a discovery source in addition to kubernetes.
 * <b>[kubewatch]</b> Kubewatch is deprecated in favor of watt.
 * <b>[lib]</b> Added supervisor package.