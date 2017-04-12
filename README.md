# Kubernetes RBACQ

[![Build Status](https://travis-ci.org/sbueringer/kubernetes-rbacq.svg?branch=master)](https://travis-ci.org/sbueringer/kubernetes-rbacq)

RBACQ simplifies querying Subjects and Rights specified in Kubernetes through Roles/ClusterRoles and RoleBindings/ClusterRoleBindings.

# Installation

## Binary

Go to the [releases](https://github.com/sbueringer/kubernetes-rbacq/releases) page and download the Linux or Windows version. Put the binary to somewhere you want (on UNIX-y systems, /usr/local/bin or the like). Make sure it has execution bits turned on.

# Basic Usage

RBACQ is build with [Cobra](https://github.com/spf13/cobra) so the CLI is build in a familiar way (Cobra is also used in Docker and Kubernetes).

To print a description what RBACQ can do, just execute:
```bash
$ ./rbacq
rbacq simplifies querying the Kubernetes RBAC API

Usage:
  rbacq [command]

Available Commands:
  get         Displays one or many resources
  help        Help about any command

Flags:
  -a, --all-namespaces      Specifies that all Namespaces should be queried (default "false")
  -c, --cluster-wide        Search cluster-wide (which includes ClusterRoles & ClusterRolebindings)
  -k, --kubeconfig string   Path to the kubeconfig file to use for CLI requests (default "$HOME\\.kube\\config")
  -n, --namespace string    Specifies the Namespace in which to query (default "default")
  -s, --system              Show also System Objects (default "false")

Use "rbacq [command] --help" for more information about a command.
```

To further explore the CLI execute the following: (and so on)
```bash
$ ./rbacq get
You must specify the type of resource to get. Valid resource types are:

        * subjects (aka 'sub')
        * rights (aka 'r')
```

```bash
$ ./rbacq get subjects --help
Displays one or many resources

Usage:
  rbacq get [RESOURCE-TYPE] [flags]

Flags:
  -o, --output string   Set jsonpath e.g. with -o jsonpath='{.kind}:{.Name}'

Global Flags:
  -a, --all-namespaces      Specifies that all Namespaces should be queried (default "false")
  -c, --cluster-wide        Search cluster-wide (which includes ClusterRoles & ClusterRolebindings)
  -k, --kubeconfig string   Path to the kubeconfig file to use for CLI requests (default "C:\\Users\\SBUERIN\\.kube\\config")
  -n, --namespace string    Specifies the Namespace in which to query (default "default")
  -s, --system              Show also System Objects (default "false")

```

## Subjects

Subjects used in RoleBindings can be queried with `./rbaq get subjects`. The subjects are queried per default in the default Namespace. The following flags can modify this behaviour:
* -n \<namespace\>: search in a specific Namespace 
* -a: search in all Namespaces 
* -c: search cluster-wide, which means that also ClusterRoles & ClusterRoleBindings are queried 
* -s: also show System objects 

Examples:

List Subjects in kube-system (including System objects):
```bash
$ ./rbacq -n kube-system get subjects -s
Subjects defined in RoleBindings
    Namespace: kube-system
        ServiceAccount:kube-system:token-cleaner
            Role: system:controller:token-cleaner
                 secrets: [delete get list watch]
                 events: [create patch update]
        ServiceAccount:infra:vault
            Role: vault:serviceaccount
                 secrets: [delete create list update]
        ServiceAccount:kube-system:bootstrap-signer
            Role: system:controller:bootstrap-signer
                 secrets: [get list watch]
```

List all Subjects matching the RegExp `.*kube-system.*` (including System objects):
```bash
$ ./rbacq -n kube-system get subjects -s .*kube-system.*
Subjects defined in RoleBindings
    Namespace: kube-system
        ServiceAccount:kube-system:token-cleaner
            Role: system:controller:token-cleaner
                 secrets: [delete get list watch]
                 events: [create patch update]
        ServiceAccount:kube-system:bootstrap-signer
            Role: system:controller:bootstrap-signer
                 secrets: [get list watch]
```

## Rights

Rights used by Roles can be queried with `./rbacq get rights`. The rights are queried per default in the default Namespace. The same flags as with Subjects can modify this behaviour.

Example:

Get all Rights from Namespace kube-system (including System):
```bash
$ ./rbacq -n kube-system get rights -s 
Rights defined in Roles
    events:
        [create patch update]: [ServiceAccount:kube-system:token-cleaner]
    secrets:
        [delete create list update]: [ServiceAccount:infra:i3-vault]
        [delete get list watch]: [ServiceAccount:kube-system:token-cleaner]
        [get list watch]: [ServiceAccount:kube-system:bootstrap-signer]
```

Get all Rights from Roles in default Namespaces and ClusterRoles that match `namespaces.*get` (including System):
```bash
$ ./rbacq get rights -s -c namespaces.*get
Rights defined in ClusterRoles & Roles
    namespaces: [delete get list watch]: [ServiceAccount:kube-system:namespace-controller]
    namespaces: [get]: [User:system:kube-controller-manager]

```
TODO https://github.com/dwyl/repo-badges
