/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	bookinfov1beta1 "com/bookinfocrd/api/v1beta1"
)

// BookinfoReconciler reconciles a Bookinfo object
type BookinfoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=bookinfo.com,resources=bookinfoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bookinfo.com,resources=bookinfoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=bookinfo.com,resources=bookinfoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bookinfo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *BookinfoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log.Log.Info("trigger bookinfo app reconcile!")
	var bookinfo bookinfov1beta1.Bookinfo
	if err := r.Client.Get(ctx, req.NamespacedName, &bookinfo); err != nil {
		//资源没有或k8s异常
		if apierrors.IsNotFound(err) {
			//说明执行了删除bookinfo资源
			//需要去删除（deployment+service+ingress）资源
			bookinfo.Name = req.Name
			bookinfo.Namespace = req.Namespace
			if err := r.deleteBookinfo(ctx, bookinfo); err != nil {
				msg := fmt.Sprintf("delete bookinfo app %s err %v", req.Name, err)
				log.Log.Error(err, msg)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, err
		}
	} else {
		//创建 bookinfo资源
		//同步（deployment+service+ingress）资源
		if err := r.createBookinfo(ctx, bookinfo); err != nil {
			msg := fmt.Sprintf("create bookinfo app %s err %v", req.Name, err)
			log.Log.Error(err, msg)
			bookinfo.Status.Message = msg
			r.Client.Status().Update(ctx, &bookinfo)
			return ctrl.Result{}, err
		}
		msg := fmt.Sprintf("create bookinfo app %s succcessfully", req.Name)
		bookinfo.Status.Message = msg
		r.Client.Status().Update(ctx, &bookinfo)
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookinfoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bookinfov1beta1.Bookinfo{}).
		Complete(r)
}
