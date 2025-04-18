/*
Copyright 2023 The OpenYurt Authors.

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

package v1beta1

import (
	"strings"

	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/openyurtio/openyurt/pkg/apis/apps"
	"github.com/openyurtio/openyurt/pkg/apis/apps/v1beta2"
)

func (src *NodePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.NodePool)

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Type = v1beta2.NodePoolType(src.Spec.Type)
	dst.Spec.Labels = src.Spec.Labels
	dst.Spec.Annotations = src.Spec.Annotations
	dst.Spec.Taints = src.Spec.Taints
	if strings.EqualFold(src.Annotations[apps.NodePoolHostNetworkLabel], "true") {
		dst.Spec.HostNetwork = true
	}

	dst.Status.ReadyNodeNum = src.Status.ReadyNodeNum
	dst.Status.UnreadyNodeNum = src.Status.UnreadyNodeNum
	dst.Status.Nodes = src.Status.Nodes

	// Set interconnectivity to false which will not use leader election strategy or reuse list/watch events
	dst.Spec.InterConnectivity = false
	dst.Spec.LeaderElectionStrategy = string(v1beta2.ElectionStrategyRandom)
	dst.Spec.LeaderReplicas = 1

	klog.V(4).Infof("convert from v1beta to v1beta2 for nodepool %s", dst.Name)

	return nil
}

func (dst *NodePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.NodePool)

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Type = NodePoolType(src.Spec.Type)
	dst.Spec.Labels = src.Spec.Labels
	dst.Spec.Annotations = src.Spec.Annotations
	dst.Spec.Taints = src.Spec.Taints

	dst.Status.ReadyNodeNum = src.Status.ReadyNodeNum
	dst.Status.UnreadyNodeNum = src.Status.UnreadyNodeNum
	dst.Status.Nodes = src.Status.Nodes

	if src.Spec.HostNetwork {
		if dst.Annotations == nil {
			dst.Annotations = make(map[string]string)
		}
		dst.Annotations[apps.NodePoolHostNetworkLabel] = "true"
	}

	klog.V(4).Infof("convert from v1beta2 to v1beta1 for nodepool %s", dst.Name)
	return nil
}
