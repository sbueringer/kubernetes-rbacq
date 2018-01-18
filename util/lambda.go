// Copyright © 2018 Stefan Bueringer
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import "k8s.io/client-go/pkg/apis/rbac/v1beta1"

//Returns the first index of the target string t, or -1 if no match is found.
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Returns true if the target string t is in the slice.
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

//Returns true if one of the strings in the slice satisfies the predicate f.
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

//Returns true if all of the strings in the slice satisfy the predicate f.
func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

//Returns a new slice containing all strings in the slice that satisfy the predicate f.
func SubjectFilter(vs []v1beta1.Subject, f func(subject v1beta1.Subject) bool) []v1beta1.Subject {
	vsf := make([]v1beta1.Subject, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

//Returns a new slice containing all strings in the slice that satisfy the predicate f.
func RoleFilter(vs []v1beta1.Role, f func(subject v1beta1.Role) bool) []v1beta1.Role {
	vsf := make([]v1beta1.Role, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

//Returns a new slice containing all strings in the slice that satisfy the predicate f.
func ClusterRoleFilter(vs []v1beta1.ClusterRole, f func(subject v1beta1.ClusterRole) bool) []v1beta1.ClusterRole {
	vsf := make([]v1beta1.ClusterRole, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

//Returns a new slice containing all strings in the slice that satisfy the predicate f.
func RoleBindingFilter(vs []v1beta1.RoleBinding, f func(subject v1beta1.RoleBinding) bool) []v1beta1.RoleBinding {
	vsf := make([]v1beta1.RoleBinding, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

//Returns a new slice containing all strings in the slice that satisfy the predicate f.
func ClusterRoleBindingFilter(vs []v1beta1.ClusterRoleBinding, f func(subject v1beta1.ClusterRoleBinding) bool) []v1beta1.ClusterRoleBinding {
	vsf := make([]v1beta1.ClusterRoleBinding, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

//Returns a new slice containing all strings in the slice that satisfy the predicate f.
func SubjectRoleRefFilter(vs map[v1beta1.Subject][]v1beta1.RoleRef, f func(subject v1beta1.Subject) bool) map[v1beta1.Subject][]v1beta1.RoleRef {
	vsf := make(map[v1beta1.Subject][]v1beta1.RoleRef, 0)
	for k, v := range vs {
		if f(k) {
			vsf[k] = v
		}
	}
	return vsf
}

//Returns a new slice containing the results of applying the function f to each string in the original slice.
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
