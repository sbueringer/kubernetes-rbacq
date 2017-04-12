package query

import (
	"os"
	"strings"
	"regexp"
	"fmt"
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/rbac/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/jsonpath"

	"github.com/sbueringer/kubernetes-rbacq/logger"
	"github.com/sbueringer/kubernetes-rbacq/util"
)

var KubeCfgFile string
var Namespace string
var AllNamespaces bool
var System bool
var ClusterWide bool
var Output string
var jsonPathSet bool

var clientset *kubernetes.Clientset

var clusterRoleList *v1beta1.ClusterRoleList
var clusterRoleBindingList *v1beta1.ClusterRoleBindingList
var roleList *v1beta1.RoleList
var roleBindingList *v1beta1.RoleBindingList

func InitKubeCfg() {
	var config *rest.Config

	config, err := clientcmd.BuildConfigFromFlags("", KubeCfgFile)
	logger.HandleError(err)

	clientset, err = kubernetes.NewForConfig(config)
	logger.HandleError(err)

	if Output != "" && strings.HasPrefix(Output, "jsonpath=") {
		jsonPathSet = true
	}
}

func GetRights(args []string) {
	var rightsFilter *regexp.Regexp
	if len(args) > 1 {
		rightsFilter = regexp.MustCompile(args[1])
	}

	clusterRoleList, err := clientset.RbacV1beta1Client.ClusterRoles().List(v1.ListOptions{})
	logger.HandleError(err)

	if AllNamespaces {
		roleList, err = clientset.RbacV1beta1Client.Roles("").List(v1.ListOptions{})
		logger.HandleError(err)
	} else {
		roleList, err = clientset.RbacV1beta1Client.Roles(Namespace).List(v1.ListOptions{})
		logger.HandleError(err)
	}

	clusterRoleBindingList, err = clientset.RbacV1beta1Client.ClusterRoleBindings().List(v1.ListOptions{})
	logger.HandleError(err)

	if AllNamespaces {
		roleBindingList, err = clientset.RbacV1beta1Client.RoleBindings("").List(v1.ListOptions{})
		logger.HandleError(err)
	} else {
		roleBindingList, err = clientset.RbacV1beta1Client.RoleBindings(Namespace).List(v1.ListOptions{})
		logger.HandleError(err)
	}

	clusterRoles := clusterRoleList.Items
	if !System {
		// filter for System (filter System Subjects)
		clusterRoles = util.ClusterRoleFilter(clusterRoles, func(c v1beta1.ClusterRole) bool { return !strings.HasPrefix(c.Name, "system:") })
	}

	roles := roleList.Items
	if !System {
		// filter for System (filter System Subjects)
		roles = util.RoleFilter(roles, func(r v1beta1.Role) bool { return !strings.HasPrefix(r.Name, "system:") })
	}

	if !jsonPathSet {
		if ClusterWide {
			logger.Return.Println("Rights defined in ClusterRoles & Roles")
		} else {
			logger.Return.Println("Rights defined in Roles")
		}
	}

	var policyRuleSubjectMap map[string][]v1beta1.Subject = make(map[string][]v1beta1.Subject)
	var resourceKeyMap map[string][]string = make(map[string][]string)
	if ClusterWide {
		for _, clusterRole := range clusterRoles {
			subjects := getSubjectsForClusterRole(clusterRole)
			if len(subjects) > 0 {
				for _, policyRule := range clusterRole.Rules {
					for _, resource := range policyRule.Resources {
						addPolicyRuleSubjectToMap(&resourceKeyMap, &policyRuleSubjectMap, resource, policyRule.Verbs, subjects)
					}
					for _, nonResourceURLs := range policyRule.NonResourceURLs {
						addPolicyRuleSubjectToMap(&resourceKeyMap, &policyRuleSubjectMap, nonResourceURLs, policyRule.Verbs, subjects)
					}
				}
			} else {
				logger.Debug.Printf("Unmapped Roles: %s", clusterRole.Name)
			}

		}
	}
	for _, role := range roles {
		subjects := getSubjectsForRole(role)
		for _, policyRule := range role.Rules {
			if len(subjects) > 0 {
				for _, resource := range policyRule.Resources {
					addPolicyRuleSubjectToMap(&resourceKeyMap, &policyRuleSubjectMap, resource, policyRule.Verbs, subjects)
				}
				for _, nonResourceURLs := range policyRule.NonResourceURLs {
					addPolicyRuleSubjectToMap(&resourceKeyMap, &policyRuleSubjectMap, nonResourceURLs, policyRule.Verbs, subjects)
				}
			} else {
				logger.Debug.Printf("Unmapped Roles: %s", role.Name)
			}
		}
	}
	// To store the resources in slice in sorted order
	var resources []string
	for k := range resourceKeyMap {
		resources = append(resources, k)
	}
	sort.Strings(resources)

	for _, resource := range resources {
		rights := resourceKeyMap[resource]
		sort.Strings(rights)
		var output string
		for _, right := range rights {
			if rightsFilter == nil || rightsFilter.MatchString(right) {
				var subjects []string = []string{}
				for _, subject := range policyRuleSubjectMap[right] {
					subjects = append(subjects, getFullSubjectName(subject))
				}
				output += fmt.Sprintf("\t\t%s: %v\n", strings.TrimPrefix(right, resource), subjects)
			}
		}
		if output != "" {
			logger.Return.Printf("\t%s:", resource)
			logger.Return.Print(output)
		}
	}
}
func addPolicyRuleSubjectToMap(resourceKeyMap *map[string][]string, policyRuleSubjectMap *map[string][]v1beta1.Subject, resource string, verbs []string, subjects []v1beta1.Subject) {
	key := fmt.Sprintf("%s%v", resource, verbs)
	if keyArray, ok := (*resourceKeyMap)[resource]; ok {
		if !util.Contains(keyArray, key) {
			(*resourceKeyMap)[resource] = append(keyArray, key)
		}
	} else {
		(*resourceKeyMap)[resource] = []string{key}
	}

	if subjectsArray, ok := (*policyRuleSubjectMap)[key]; ok {

		(*policyRuleSubjectMap)[key] = appendSubjectsDistinct(subjectsArray, subjects)
	} else {
		(*policyRuleSubjectMap)[key] = makeSubjectsDistinct(subjects)
	}
}

func getSubjectsForClusterRole(role v1beta1.ClusterRole) []v1beta1.Subject {
	var subjects []v1beta1.Subject = []v1beta1.Subject{}
	for _, clusterRoleBinding := range clusterRoleBindingList.Items {
		if "ClusterRole" == clusterRoleBinding.RoleRef.Kind && role.Name == clusterRoleBinding.RoleRef.Name {
			subjects = append(subjects, clusterRoleBinding.Subjects...)
		}
	}
	return subjects
}
func getSubjectsForRole(role v1beta1.Role) []v1beta1.Subject {
	var subjects []v1beta1.Subject = []v1beta1.Subject{}
	for _, roleBinding := range roleBindingList.Items {
		if "Role" == roleBinding.RoleRef.Kind && role.Name == roleBinding.RoleRef.Name {
			subjects = append(subjects, roleBinding.Subjects...)
		}
	}
	return subjects
}

func GetSubjects(args []string) {
	var subjectFilter *regexp.Regexp
	if len(args) > 1 {
		subjectFilter = regexp.MustCompile(args[1])
	}

	var err error
	clusterRoleList, err = clientset.RbacV1beta1Client.ClusterRoles().List(v1.ListOptions{})
	logger.HandleError(err)

	if AllNamespaces {
		roleList, err = clientset.RbacV1beta1Client.Roles("").List(v1.ListOptions{})
		logger.HandleError(err)
	} else {
		roleList, err = clientset.RbacV1beta1Client.Roles(Namespace).List(v1.ListOptions{})
		logger.HandleError(err)
	}

	if !jsonPathSet {
		if ClusterWide {
			logger.Return.Println("Subjects defined in ClusterRolebindings & RoleBindings")
		} else {
			logger.Return.Println("Subjects defined in RoleBindings")
		}
	}

	if ClusterWide {
		var clusterSubjectRoleRefMap map[v1beta1.Subject]v1beta1.RoleRef = make(map[v1beta1.Subject]v1beta1.RoleRef)
		clusterRoleBindingList, err := clientset.RbacV1beta1Client.ClusterRoleBindings().List(v1.ListOptions{})
		logger.HandleError(err)

		clusterRoleBindings := clusterRoleBindingList.Items
		if !System {
			// filter for System (filter System Subjects)
			clusterRoleBindings = util.ClusterRoleBindingFilter(clusterRoleBindings, func(c v1beta1.ClusterRoleBinding) bool { return !strings.HasPrefix(c.Name, "system:") })
		}
		for _, clusterRoleBinding := range clusterRoleBindings {
			for _, subject := range clusterRoleBinding.Subjects {
				if subjectFilter == nil || subjectFilter.MatchString(getFullSubjectName(subject)) {
					clusterSubjectRoleRefMap[subject] = clusterRoleBinding.RoleRef
				}
			}
		}
		if !jsonPathSet {
			logger.Return.Println("\tCluster-wide:")
		}
		printSubjects(clusterSubjectRoleRefMap, "User")
		printSubjects(clusterSubjectRoleRefMap, "Group")
		printSubjects(clusterSubjectRoleRefMap, "ServiceAccount")
	}

	var namespaceSubjectRoleRefMap map[v1beta1.Subject]v1beta1.RoleRef = make(map[v1beta1.Subject]v1beta1.RoleRef)
	var roleBindingList *v1beta1.RoleBindingList
	if AllNamespaces {
		roleBindingList, err = clientset.RbacV1beta1Client.RoleBindings("").List(v1.ListOptions{})
		logger.HandleError(err)
	} else {
		roleBindingList, err = clientset.RbacV1beta1Client.RoleBindings(Namespace).List(v1.ListOptions{})
		logger.HandleError(err)
	}

	roleBindings := roleBindingList.Items
	if !System {
		// filter for System (filter System Subjects)
		roleBindings = util.RoleBindingFilter(roleBindings, func(r v1beta1.RoleBinding) bool { return !strings.HasPrefix(r.Name, "system:") })
	}
	for _, roleBinding := range roleBindings {
		for _, subject := range roleBinding.Subjects {
			if subjectFilter == nil || subjectFilter.MatchString(getFullSubjectName(subject)) {
				namespaceSubjectRoleRefMap[subject] = roleBinding.RoleRef
			}
		}
	}
	if !jsonPathSet {
		logger.Return.Printf("\tNamespace: %s", Namespace)
	}
	printSubjects(namespaceSubjectRoleRefMap, "User")
	printSubjects(namespaceSubjectRoleRefMap, "Group")
	printSubjects(namespaceSubjectRoleRefMap, "ServiceAccount")
}

func printSubjects(subjectRoleRefMap map[v1beta1.Subject]v1beta1.RoleRef, kind string) {
	var jsonPath *jsonpath.JSONPath
	if jsonPathSet {
		jsonPath = jsonpath.New("jsonpath")

		jsonPathString := strings.Split(Output, "=")[1]    // split = take only jsonpath
		jsonPathString = strings.Trim(jsonPathString, "'") // remove leading and trailing '
		jsonPath.Parse(jsonPathString)
	}
	// filter for Kind
	subjectRoleRefMap = util.SubjectRoleRefFilter(subjectRoleRefMap, func(s v1beta1.Subject) bool { return s.Kind == kind })
	for subject, roleRef := range subjectRoleRefMap {
		if jsonPathSet {
			jsonPath.Execute(os.Stdout, subject)
			os.Stdout.WriteString("\n")
		} else {
			printSubject(subject, roleRef)
		}
	}
}

func printSubject(subject v1beta1.Subject, roleRef v1beta1.RoleRef) {
	// print Subject
	logger.Return.Println("\t\t" + getFullSubjectName(subject))
	// print RoleRelf
	logger.Return.Println("\t\t\t" + roleRef.Kind + ": " + roleRef.Name)
	// print Role Details
	policyRules := getPolicyRules(&roleRef)
	if policyRules != nil {
		for _, policyRule := range policyRules {
			for _, resource := range policyRule.Resources {
				logger.Return.Printf("\t\t\t\t %s: %v", resource, policyRule.Verbs)
			}
			for _, nonResourceURLs := range policyRule.NonResourceURLs {
				logger.Return.Printf("\t\t\t\t %s: %v", nonResourceURLs, policyRule.Verbs)
			}
		}
	}
}

func getFullSubjectName(subject v1beta1.Subject) (string) {
	if subject.Namespace != "" {
		return subject.Kind + ":" + subject.Namespace + ":" + subject.Name
	}
	return subject.Kind + ":" + subject.Name
}
func getPolicyRules(roleRef *v1beta1.RoleRef) ([] v1beta1.PolicyRule) {
	switch roleRef.Kind {
	case "ClusterRole":
		for _, clusterRole := range clusterRoleList.Items {
			if clusterRole.Name == roleRef.Name {
				return clusterRole.Rules
			}
		}
	case "Role":
		for _, role := range roleList.Items {
			if role.Name == roleRef.Name {
				return role.Rules
			}
		}
	}
	return nil
}

func makeSubjectsDistinct(subjects []v1beta1.Subject) []v1beta1.Subject {
	var newSubjects []v1beta1.Subject
	for _, subject := range subjects {
		newSubjects = appendSubjectsDistinct(newSubjects , []v1beta1.Subject{subject})
	}
	return newSubjects
}

func appendSubjectsDistinct(subjects []v1beta1.Subject, subjectsToAdd []v1beta1.Subject) []v1beta1.Subject {
	for _, subjectToAdd := range subjectsToAdd {
		alreadyContained := false
		for _, subject := range subjects {
			if subjectToAdd.Kind == "ServiceAccount" {
				if subject.Kind == subjectToAdd.Kind && subject.Name == subjectToAdd.Name && subject.Namespace == subjectToAdd.Namespace {
					alreadyContained = true
				}
			} else {
				if subject.Kind == subjectToAdd.Kind && subject.Name == subjectToAdd.Name {
					alreadyContained = true
				}
			}
		}
		if !alreadyContained {
			subjects = append(subjects, subjectToAdd)
		}
	}
	return subjects
}
