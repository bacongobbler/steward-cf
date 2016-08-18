package claim

import (
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"golang.org/x/net/context"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// StartControlLoops calls StartControlLoop for each namespace in namespaces. For each call to StartControlLoop, it calls evtNamespacer.Events(ns) to create a new ConfigMapInterface. Because each StartControlLoop call is done in a new goroutine, this function need not be called in its own goroutine.
//
// In order to stop all loops, pass a cancellable context to this function and call its cancel function
func StartControlLoops(
	ctx context.Context,
	evtNamespacer InteractorNamespacer,
	cmNamespacer kcl.ConfigMapsNamespacer,
	lookup k8s.ServiceCatalogLookup,
	lifecycler mode.Lifecycler,
	namespaces []string,
	errCh chan<- error,
) {
	for _, ns := range namespaces {
		go func(ns string) {
			evtIface := evtNamespacer.Interactor(ns)
			if err := StartControlLoop(ctx, evtIface, cmNamespacer, lookup, lifecycler); err != nil {
				logger.Errorf("failed to start control loop for namespace %s, skipping (%s)", ns, err)
				errCh <- err
				return
			}
		}(ns)
	}
}