/*
Copyright 2019 The Knative Authors

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

package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	"knative.dev/eventing/pkg/apis/duck"
	"knative.dev/pkg/apis"
)

const (
	// CouchDbConditionReady has status True when the CouchDbSource is ready to send events.
	CouchDbConditionReady = apis.ConditionReady

	// CouchDbConditionSinkProvided has status True when the CouchDbSource has been configured with a sink target.
	CouchDbConditionSinkProvided apis.ConditionType = "SinkProvided"

	// CouchDbConditionDeployed has status True when the CouchDbSource has had it's deployment created.
	CouchDbConditionDeployed apis.ConditionType = "Deployed"
)

var CouchDbCondSet = apis.NewLivingConditionSet(
	CouchDbConditionSinkProvided,
	CouchDbConditionDeployed,
)

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*CouchDbSource) GetConditionSet() apis.ConditionSet {
	return CouchDbCondSet
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *CouchDbSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return CouchDbCondSet.Manage(s).GetCondition(t)
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *CouchDbSourceStatus) InitializeConditions() {
	CouchDbCondSet.Manage(s).InitializeConditions()
}

// MarkSink sets the condition that the source has a sink configured.
func (s *CouchDbSourceStatus) MarkSink(uri *apis.URL) {
	s.SinkURI = uri
	if !uri.IsEmpty() {
		CouchDbCondSet.Manage(s).MarkTrue(CouchDbConditionSinkProvided)
	} else {
		CouchDbCondSet.Manage(s).MarkUnknown(CouchDbConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.%s", "")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (s *CouchDbSourceStatus) MarkNoSink(reason, messageFormat string, messageA ...interface{}) {
	CouchDbCondSet.Manage(s).MarkFalse(CouchDbConditionSinkProvided, reason, messageFormat, messageA...)
}

// PropagateDeploymentAvailability uses the availability of the provided Deployment to determine if
// CouchDbConditionDeployed should be marked as true or false.
func (s *CouchDbSourceStatus) PropagateDeploymentAvailability(d *appsv1.Deployment) {
	if duck.DeploymentIsAvailable(&d.Status, false) {
		CouchDbCondSet.Manage(s).MarkTrue(CouchDbConditionDeployed)
	} else {
		// I don't know how to propagate the status well, so just give the name of the Deployment
		// for now.
		CouchDbCondSet.Manage(s).MarkFalse(CouchDbConditionDeployed, "DeploymentUnavailable", "The Deployment '%s' is unavailable.", d.Name)
	}
}

// IsReady returns true if the resource is ready overall.
func (s *CouchDbSourceStatus) IsReady() bool {
	return CouchDbCondSet.Manage(s).IsHappy()
}
