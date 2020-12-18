package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/cybozu-go/moco"
	mocov1alpha1 "github.com/cybozu-go/moco/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getPassword(ctx context.Context, clusterUniqueName, user string) (string, error) {
	secret := &corev1.Secret{}
	err := kubeClient.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      "root-password-" + clusterUniqueName,
	}, secret)
	if err != nil {
		return "", err
	}

	userPassKeys := map[string]string{
		"root":            moco.RootPasswordEnvName,
		moco.WritableUser: moco.WritablePasswordEnvName,
		moco.ReadOnlyUser: moco.ReadOnlyPasswordEnvName,
	}
	key, ok := userPassKeys[user]
	if !ok {
		return "", errors.New("unknown user: " + user)
	}
	password, ok := secret.Data[key]
	if !ok {
		return "", errors.New("unknown user: " + user)
	}
	return string(password), nil
}

func getPodName(ctx context.Context, cluster *mocov1alpha1.MySQLCluster) (string, error) {
	if mysqlConfig.index >= int(cluster.Spec.Replicas) {
		return "", errors.New("index should be smaller than replicas")
	}
	index := mysqlConfig.index
	if mysqlConfig.index < 0 {
		if cluster.Status.CurrentPrimaryIndex != nil {
			index = *cluster.Status.CurrentPrimaryIndex
		} else {
			return "", errors.New("primary instance not found")
		}
	}

	return fmt.Sprintf("%s-%d", moco.UniqueName(cluster), index), nil
}