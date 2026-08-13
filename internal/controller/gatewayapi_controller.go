package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	platformv1 "github.com/jimmyjoy/gateway-operator/api/v1"
)

// GatewayAPIReconciler reconciles a GatewayAPI object
type GatewayAPIReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=platform.mac.com,resources=gatewayapis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=platform.mac.com,resources=gatewayapis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=platform.mac.com,resources=gatewayapis/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *GatewayAPIReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {

	log := logf.FromContext(ctx)

	// 1. Read the GatewayAPI CR
	gatewayAPI := &platformv1.GatewayAPI{}

	err := r.Get(ctx, req.NamespacedName, gatewayAPI)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// CR was deleted, nothing left to reconcile
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	log.Info(
		"Reconciling GatewayAPI",
		"name", gatewayAPI.Name,
		"namespace", gatewayAPI.Namespace,
	)

	// 2. Desired ConfigMap name
	configMapName := gatewayAPI.Name + "-config"

	// 3. Try to read existing ConfigMap
	existingConfigMap := &corev1.ConfigMap{}

	err = r.Get(
		ctx,
		types.NamespacedName{
			Name:      configMapName,
			Namespace: gatewayAPI.Namespace,
		},
		existingConfigMap,
	)

	// 4. ConfigMap does not exist -> create it
	if apierrors.IsNotFound(err) {

		newConfigMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: gatewayAPI.Namespace,
			},

			Data: map[string]string{
				"teamName":       gatewayAPI.Spec.TeamName,
				"apiName":        gatewayAPI.Spec.APIName,
				"authentication": gatewayAPI.Spec.Authentication.Type,
				"backendURL":     gatewayAPI.Spec.Backend.URL,
			},
		}

		// Make GatewayAPI the owner of the ConfigMap
		if err := controllerutil.SetControllerReference(
			gatewayAPI,
			newConfigMap,
			r.Scheme,
		); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, newConfigMap); err != nil {
			return ctrl.Result{}, err
		}

		log.Info(
			"Created ConfigMap",
			"name", configMapName,
		)

		return ctrl.Result{}, nil
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	// 5. ConfigMap already exists -> ensure it matches desired state
	existingConfigMap.Data = map[string]string{
		"teamName":       gatewayAPI.Spec.TeamName,
		"apiName":        gatewayAPI.Spec.APIName,
		"authentication": gatewayAPI.Spec.Authentication.Type,
		"backendURL":     gatewayAPI.Spec.Backend.URL,
	}

	if err := r.Update(ctx, existingConfigMap); err != nil {
		return ctrl.Result{}, err
	}

	log.Info(
		"ConfigMap reconciled",
		"name", configMapName,
	)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayAPIReconciler) SetupWithManager(
	mgr ctrl.Manager,
) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1.GatewayAPI{}).
		Owns(&corev1.ConfigMap{}).
		Named("gatewayapi").
		Complete(r)
}
