# longroutine
Long running go routines (if you absolutely need them) for Kubernetes controllers

## Usage

Add `longroutine.Starter` as one of your reconciler fields. Start your go routines with StartSingle() method:

```go
package whatever

import (
	"github.com/dispatchframework/longroutine"
)

type reconciler struct {
	RoutineStarter longroutine.SingleStarter
}

func NewReconciler() *reconciler {
	return &reconciler{
		RoutineStarter: longroutine.NewSingleStarter(),
	}
}

func (r *reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// skip ...	

	// starting a long running go routine, e.g. a ClusterAPI cluster upgrade
	// if a routine is already running for this upgradeID, the new one will not be started
	r.RoutineStarter.StartSingle(upgradeID, func() {
		upgrader.Upgrade()
	})

	// skip ...
	
	return reconcile.Result{}, nil
}
```
