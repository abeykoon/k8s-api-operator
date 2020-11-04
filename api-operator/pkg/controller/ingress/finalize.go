package ingress

import "github.com/wso2/k8s-api-operator/api-operator/pkg/controller/common"

const (
	// finalizerName represents the name of ingress finalizer handled by this controller
	finalizerName = "wso2.microgateway/ingress.finalizer"
)

func finalizeDeletion(requestInfo *common.RequestInfo) error {
	// handle deletion with finalizers to avoid missing ingress configurations deleted while
	// restating controller, or deleted before starting controller.
	//
	// Ingress deletion delta change also handled in the update delta change flow and
	// skipping handling deletion here
	return nil
}
