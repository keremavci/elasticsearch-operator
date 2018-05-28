package k8shandler

import (
	"bytes"
	"fmt"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	v1alpha1 "github.com/t0ffel/elasticsearch-operator/pkg/apis/elasticsearch/v1alpha1"
	//"github.com/sirupsen/logrus"
)

func CreateOrUpdateConfigMaps(dpl *v1alpha1.Elasticsearch) error {
	elasticsearchCMName := dpl.Name
	owner := asOwner(dpl)

	// TODO: take all vars from CRD
	pathData := "- /elasticsearch/persistent/"
	err := createOrUpdateConfigMap(elasticsearchCMName, dpl.Namespace, dpl.Name, defaultKibanaIndexMode, pathData, false, owner)
	if err != nil {
		return fmt.Errorf("Failure creating ConfigMap %v", err)
	}
	return nil
}

func createOrUpdateConfigMap(configMapName, namespace, clusterName, kibanaIndexMode, pathData string, allowClusterReader bool, owner metav1.OwnerReference) error {
	elasticsearchCM, err := createConfigMap(configMapName, namespace, clusterName, kibanaIndexMode, pathData, allowClusterReader)
	if err != nil {
		return err
	}
	addOwnerRefToObject(elasticsearchCM, owner)
	err = action.Create(elasticsearchCM)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("Failure constructing Elasticsearch ConfigMap: %v", err)
	} else if errors.IsAlreadyExists(err) {
		// Get existing configMap to check if it is same as what we want
		existingCM := configMap(configMapName, namespace)
		err = query.Get(existingCM)
		if err != nil {
			return fmt.Errorf("Unable to get Elasticsearch cluster configMap: %v", err)
		}

		// TODO: Compare existing configMap labels, selectors and port
	}
	return nil
}

func createConfigMap(configMapName string, namespace string, clusterName string, kibanaIndexMode string, pathData string, allowClusterReader bool) (*v1.ConfigMap, error) {
	cm := configMap(configMapName, namespace)
	cm.Data = map[string]string{}
	buf := &bytes.Buffer{}
	if err := renderEsYml(buf, allowClusterReader, kibanaIndexMode, pathData); err != nil {
		return cm, err
	}
	cm.Data["elasticsearch.yml"] = buf.String()

	buf = &bytes.Buffer{}
	if err := renderLog4j2Properties(buf, defaultRootLogger); err != nil {
		return cm, err
	}
	cm.Data["log4j2.properties"] = buf.String()

	return cm, nil
}

// configMap returns a v1.ConfigMap object
func configMap(configMapName string, namespace string) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
	}
}
