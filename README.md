# Kubernetes RBACQ

[![Build Status](https://travis-ci.org/sbueringer/kubernetes-rbacq.svg?branch=master)](https://travis-ci.org/sbueringer/kubernetes-rbacq)

Query Tool for Kubernetes RBAC Objects

# Query

 RBAC Window bauen zum erklären / debuggen / mit query welche roles wer hat welche roles welche Rechte hat...

# Query Use Cases

Data:
* Subjects, RoleBindings, Roles, Rights

Use Cases

Rights:
* who has which right:
    * Roles
        * Subjects
    *Subjects

Subject:
* Which Subjects are binded
    * which roles has a subject
    * which rights has a subject

# Api

rbacq get subjects 
(with Roles & Rights)

rbac get subject <regexp>

TODO verify functionality, create documentation

# Check Rights for user / current user


kubectl create --v=8 -f -  << __EOF__
{
  "apiVersion": "authorization.k8s.io/v1",
  "kind": "SubjectAccessReview",
  "spec": {
    "resourceAttributes": {
      "group": "",
      "verb": "get",
      "resource": "endpoints"
    },
    "user": "system:kube-controller-manager"
  }
}
__EOF__
