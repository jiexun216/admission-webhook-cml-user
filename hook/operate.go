package hook

import (
	"encoding/json"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
)

//1.# console
//monitoring-platform-access
//2.# cml
//associatedCRP
//4.# warehouse
//istio-injection
//5.# compute
//istio-injection
//6.# implala
//istio-injection
//7.# monitoring
//cdp.cloudera/version

//3.# cml-user
//ds-parent-namespace

//就是只有# cml-user
//ds-parent-namespace  是admission执行加resouce limit 和加securityContext；
//其他是admission执行删除resources和加securityContext
func createSpecSecurityAndModifyResourcesContextPatch(pod corev1.Pod, availableAnnotations map[string]string, annotations map[string]string, availableLabels map[string]string, labels map[string]string) ([]byte, error) {
	var patch []patchOperation
	// update Annotation to set admissionWebhookAnnotationStatusKey: "mutated"
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)
	// add labels
	patch = append(patch, updateLabels(availableLabels, labels)...)

	// update pod spec securityContext
	replaceSecurityContext := patchOperation{
		Op:   "replace",
		Path: "/spec/securityContext",
		Value: &corev1.PodSecurityContext{
			RunAsUser:  &securityContextValue,
			RunAsGroup: &securityContextValue,
			FSGroup:    &securityContextValue,
		},
	}
	glog.Infof("modify  pod spec securityContext for value: %v", replaceSecurityContext)
	patch = append(patch, replaceSecurityContext)

	// update pod container resource
	replaceResourse := patchOperation{
		Op:    "replace",
		Path:  "/spec/containers/0/resources",
		Value: assembleResourceRequirements(pod),
	}
	glog.Infof("update  pod container resource for value: %v", replaceResourse)
	patch = append(patch, replaceResourse)

	return json.Marshal(patch)
}




func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{}
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}

func updateLabels(target map[string]string, added map[string]string) (patch []patchOperation) {
	values := make(map[string]string)
	for key, value := range added {
		if target == nil || target[key] == "" {
			values[key] = value
		}
	}
	patch = append(patch, patchOperation{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: values,
	})
	return patch
}

func assembleResourceRequirements(pod corev1.Pod) corev1.ResourceRequirements {
	res := corev1.ResourceRequirements{}
	// jiexun default requests have value, put the requests value to the limits value
	res.Limits = pod.Spec.Containers[0].Resources.Requests
	res.Requests = pod.Spec.Containers[0].Resources.Requests
	return res
}
