// Copyright 2023 Upbound Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/upbound/up-sdk-go/apis/common"
)

const (
	// ResourceCredentialsSecretInClusterKubeconfigKey is the key in the
	// connection secret of the ControlPlane that contains the kubeconfig
	// to be used by running pods in the cluster.
	ResourceCredentialsSecretInClusterKubeconfigKey = "kubeconfig-incluster"

	// ConditionMessageAnnotationKey is the key for the message shown in the
	// message column in kubectl.
	ConditionMessageAnnotationKey = "internal.spaces.upbound.io/message"

	// ControlPlaneGroupLabelKey is the key used to identify namespaces as groups. The
	// value will be "true".
	ControlPlaneGroupLabelKey = "spaces.upbound.io/group"

	// ControlPlaneGroupProtectionKey is the key used to prevent deletion of groups
	// via the Spaces API. Deletion will not be protected if the underlying namespace
	// is deleted.
	ControlPlaneGroupProtectionKey = "spaces.upbound.io/group-deletion-protection"

	// The value expected by ControlPlaneGroupLabelKey and ManagedBySpacesAPILabelKey.
	LabelValueTrue = "true"
)

const (
	// FeaturesAnnotation is an optional annotation that enables features
	// gates within the control plane compositions. Value should be defined
	// as an inline map of key value pairs expressing features to be enabled.
	// For example: `{"featureA": true,"featureB": false}`. The default value is
	// empty (no features enabled).
	FeaturesAnnotation = "internal.spaces.upbound.io/features"
	// TierLimitsAnnotation is an optional annotation that specifies the limits
	// applied to the control plane, as metered by the mxp-account-gate. These
	// limits are only applicable when the account gate is enabled using the
	// features annotation.
	TierLimitsAnnotation = "internal.spaces.upbound.io/tier-limits"
	// ManagedBySpacesAPILabelKey is the key used to identify resources managed through spaces-api. The
	// value will be "true".
	ManagedBySpacesAPILabelKey = "internal.spaces.upbound.io/managed-by-spaces-api"
)

// CrossplaneUpgradeChannel is the channel for Crossplane upgrades.
type CrossplaneUpgradeChannel string

const (
	// CrossplaneUpgradeNone disables auto-upgrades and keeps the control plane at its current version of Crossplane.
	CrossplaneUpgradeNone CrossplaneUpgradeChannel = "None"

	// CrossplaneUpgradePatch automatically upgrades the control plane to the latest supported patch version when it
	// becomes available while keeping the minor version the same.
	CrossplaneUpgradePatch CrossplaneUpgradeChannel = "Patch"

	// CrossplaneUpgradeStable automatically upgrades the control plane to the latest supported patch release on minor
	// version N-1, where N is the latest supported minor version.
	CrossplaneUpgradeStable CrossplaneUpgradeChannel = "Stable"

	// CrossplaneUpgradeRapid automatically upgrades the cluster to the latest supported patch release on the latest
	// supported minor version.
	CrossplaneUpgradeRapid CrossplaneUpgradeChannel = "Rapid"
)

// CrossplaneState is the running state for the crossplane and provider workloads.
type CrossplaneState string

const (
	// CrossplaneStateRunning switches the crossplane and provider workloads to
	// the running state by scaling up them.
	CrossplaneStateRunning CrossplaneState = "Running"

	// CrossplaneStatePaused switches the crossplane and provider workloads to
	// the paused state by scaling down them.
	CrossplaneStatePaused CrossplaneState = "Paused"
)

const (
	// ClassDefault is the default class for the control plane.
	ClassDefault = "default"
)

// ManagedResourceActivationMode defines which ManagedResourceDefinitions a control
// plane activates by default.
type ManagedResourceActivationMode string

const (
	// ManagedResourceActivationModeAll activates every
	// ManagedResourceDefinition offered by the installed providers. This is
	// Crossplane's own default.
	ManagedResourceActivationModeAll ManagedResourceActivationMode = "All"

	// ManagedResourceActivationModeNone activates no
	// ManagedResourceDefinition, letting providers scale to zero until
	// something activates the managed resources actually in use.
	ManagedResourceActivationModeNone ManagedResourceActivationMode = "None"
)

// CrossplaneDefaultManagedResourceActivationPolicySpec configures the
// ManagedResourceActivationPolicy named "default" in a control plane. It does
// not configure any other activation policy: a control plane's owner remains
// free to create their own. It only applies to control
// planes running Crossplane 2.x, where ManagedResourceDefinitions and
// activation policies exist.
type CrossplaneDefaultManagedResourceActivationPolicySpec struct {
	// Mode selects which ManagedResourceDefinitions the "default"
	// ManagedResourceActivationPolicy activates.
	// - All: activate every ManagedResourceDefinition the installed providers
	//   offer. This is Crossplane's own default.
	// - None: activate none of them, so providers can scale to zero until a
	//   ManagedResourceActivationPolicy activates what is actually in use.
	//
	// Activation is one-way in Crossplane: a ManagedResourceDefinition that is
	// already Active cannot go back to Inactive. None therefore prevents
	// further activation rather than undoing it, and only delivers provider
	// scale-to-zero on a control plane whose managed resources were never
	// activated.
	//
	// Whenever this is set, Spaces reconciles the "default"
	// ManagedResourceActivationPolicy to match, so changes take effect on a
	// running control plane and manual edits to that object are reverted.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=All;None
	Mode ManagedResourceActivationMode `json:"mode"`
}

// CrossplaneAutoUpgradeSpec defines the auto upgrade policy for Crossplane.
type CrossplaneAutoUpgradeSpec struct {
	// Channel defines the upgrade channels for Crossplane. We support the following channels where 'Stable' is the
	// default:
	// - None: disables auto-upgrades and keeps the control plane at its current version of Crossplane.
	// - Patch: automatically upgrades the control plane to the latest supported patch version when it
	//   becomes available while keeping the minor version the same.
	// - Stable: automatically upgrades the control plane to the latest supported patch release on minor
	//   version N-1, where N is the latest supported minor version.
	// - Rapid: automatically upgrades the cluster to the latest supported patch release on the latest
	//   supported minor version.
	// +optional
	// +kubebuilder:default="Stable"
	// +kubebuilder:validation:Enum="None";"Patch";"Stable";"Rapid"
	Channel *CrossplaneUpgradeChannel `json:"channel,omitempty"`
}

// CrossplaneSpec defines the configuration for Crossplane.
type CrossplaneSpec struct {
	// Version is the version of Universal Crossplane to install.
	// +optional
	// +kubebuilder:validation:XValidation:rule=(self.matches('^[^v].*')),message=The version must not start with a leading 'v'
	Version *string `json:"version,omitempty"`

	// AutoUpgrades defines the auto upgrade configuration for Crossplane.
	// +optional
	// +kubebuilder:default={"channel":"Stable"}
	AutoUpgradeSpec *CrossplaneAutoUpgradeSpec `json:"autoUpgrade,omitempty"`

	// State defines the state for crossplane and provider workloads. We support
	// the following states where 'Running' is the default:
	// - Running: Starts/Scales up all crossplane and provider workloads in the ControlPlane
	// - Paused: Pauses/Scales down all crossplane and provider workloads in the ControlPlane
	// +optional
	// +kubebuilder:validation:Enum=Running;Paused
	// +kubebuilder:default=Running
	State *CrossplaneState `json:"state,omitempty"`

	// DefaultManagedResourceActivationPolicy configures the "default"
	// ManagedResourceActivationPolicy of the control plane, on Crossplane 2.x
	// only. It overrides the Space-wide default. When unset, the Space-wide
	// default applies, and when nothing is set anywhere the policy activates
	// every ManagedResourceDefinition, as Crossplane does by default. Ignored on
	// Crossplane 1.x control planes, which have no activation policies.
	// +optional
	DefaultManagedResourceActivationPolicy *CrossplaneDefaultManagedResourceActivationPolicySpec `json:"defaultManagedResourceActivationPolicy,omitempty"`
}

// A SecretReference is a reference to a secret in an arbitrary namespace.
type SecretReference struct {
	// Name of the secret.
	Name string `json:"name"`

	// Namespace of the secret. If omitted, it is equal to
	// the namespace of the resource containing this reference as a field.
	// +optional
	Namespace string `json:"namespace"`
}

// A ControlPlaneSpec represents the desired state of the ControlPlane.
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.restore) || has(self.restore)",message="restore source can not be unset"
// +kubebuilder:validation:XValidation:rule="has(oldSelf.restore) || !has(self.restore)",message="restore source can not be set after creation"
// +kubebuilder:validation:XValidation:rule="!has(self.crossplane.autoUpgrade) || self.crossplane.autoUpgrade.channel != \"None\" || self.crossplane.version != \"\"",message="\"version\" cannot be empty when upgrade channel is \"None\""
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.hostNamespace) && !has(self.hostNamespace)) || (has(oldSelf.hostNamespace) && has(self.hostNamespace) && oldSelf.hostNamespace.name == self.hostNamespace.name)",message="[[GATE:EnableHostNamespace]] hostNamespace is immutable after creation"
type ControlPlaneSpec struct {
	// WriteConnectionSecretToReference specifies the namespace and name of a
	// Secret to which any connection details for this managed resource should
	// be written. Connection details frequently include the endpoint, username,
	// and password required to connect to the managed resource.
	// This field is planned to be replaced in a future release in favor of
	// PublishConnectionDetailsTo. Currently, both could be set independently
	// and connection details would be published to both without affecting
	// each other.
	//
	// If omitted, it is defaulted to the namespace of the ControlPlane.
	// Deprecated: Use Hub or Upbound identities instead.
	// +optional
	WriteConnectionSecretToReference *SecretReference `json:"writeConnectionSecretToRef,omitempty"`

	// Crossplane defines the configuration for Crossplane.
	Crossplane CrossplaneSpec `json:"crossplane,omitempty"`

	// Restore specifies details about the control planes restore configuration.
	// +optional
	// +kubebuilder:validation:XValidation:rule="!has(oldSelf.finishedAt) || oldSelf.finishedAt == self.finishedAt",message="finishedAt is immutable once set"
	Restore *Restore `json:"restore,omitempty"`

	// [[GATE:EnableControlPlaneClasses]]
	// Class specifies the class of the control plane. This affects the
	// control plane sizing, including component replicas and resource
	// requirements. There are multiple predefined classes, with "default"
	// being the standard Spaces control plane without any additional class
	// configuration. Check the Upbound Cloud documentation for a list of all
	// available classes. Defaults to the configured default class, if not set.
	// Can only be updated to a compatible class, having the same declared type.
	// +optional
	Class string `json:"class,omitempty"`

	// [[GATE:EnableHostNamespace]]
	// HostNamespace configures the host cluster namespace for this control plane.
	// When set, Spaces provisions the control plane into the namespace specified
	// by hostNamespace.name. If the namespace does not exist, Spaces creates it.
	// hostNamespace.name is immutable after creation.
	// +optional
	HostNamespace *HostNamespaceSpec `json:"hostNamespace,omitempty"`
}

// HostNamespaceDeletionPolicy specifies what happens to the host namespace when
// the ControlPlane is deleted.
type HostNamespaceDeletionPolicy string

const (
	// HostNamespaceDeletionPolicyDelete deletes the namespace and all resources
	// inside it when the ControlPlane is deleted.
	HostNamespaceDeletionPolicyDelete HostNamespaceDeletionPolicy = "Delete"
	// HostNamespaceDeletionPolicyOrphan preserves the namespace and all resources
	// inside it when the ControlPlane is deleted. Cluster-scoped resources
	// created by Spaces for this control plane are still removed.
	HostNamespaceDeletionPolicyOrphan HostNamespaceDeletionPolicy = "Orphan"
)

// HostNamespaceSpec configures the host cluster namespace for a control plane.
type HostNamespaceSpec struct {
	// Name is the name of the host cluster namespace. If the namespace does not
	// exist, Spaces creates it. Immutable after creation.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern="^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
	Name string `json:"name"`

	// DeletionPolicy specifies what happens to the host namespace when the
	// ControlPlane is deleted. Delete (the default) deletes the namespace and
	// all resources inside it. Orphan preserves the namespace and leaves all
	// resources in place; cluster-scoped resources created by Spaces are still
	// removed.
	// +optional
	// +kubebuilder:default=Delete
	// +kubebuilder:validation:Enum=Delete;Orphan
	DeletionPolicy HostNamespaceDeletionPolicy `json:"deletionPolicy,omitempty"`
}

// Restore specifies details about the backup to restore from.
type Restore struct {
	// Source of the Backup or BackupSchedule to restore from.
	// Require "restore" permission on the referenced Backup or BackupSchedule.
	// ApiGroup is optional and defaults to "spaces.upbound.io".
	// Kind is required, and the only supported kinds are Backup and
	// BackupSchedule at the moment.
	// Name is required.
	// +kubebuilder:validation:XValidation:rule="(!has(self.apiGroup) || self.apiGroup == 'spaces.upbound.io') && (self.kind == 'Backup' || self.kind == 'BackupSchedule')",message="source must be a reference to a Backup or BackupSchedule (v1alpha1)"
	// +kubebuilder:validation:XValidation:rule="oldSelf == self",message="source is immutable"
	Source common.TypedLocalObjectReference `json:"source"`

	// FinishedAt is the time at which the control plane was restored, it's not
	// meant to be set by the user, but rather by the system when the control
	// plane is restored.
	FinishedAt *metav1.Time `json:"finishedAt,omitempty"`
}

// A ControlPlaneStatus represents the observed state of a ControlPlane.
type ControlPlaneStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`
	// Message is a human-readable message indicating details about why the
	// ControlPlane is in this condition.
	Message        string `json:"message,omitempty"`
	ControlPlaneID string `json:"controlPlaneID,omitempty"`
	// FirstAvailableAt is the time at which the control plane was available for the first time.
	// +optional
	FirstAvailableAt *metav1.Time `json:"firstAvailableAt,omitempty"`
	// [[GATE:EnableControlPlaneClasses]]
	// Class holds the status about the control plane class.
	Class *ControlPlaneClassStatus `json:"class,omitempty"`
}

// ControlPlaneClassStatus represents the status of the ControlPlane's class.
type ControlPlaneClassStatus struct {
	// Current holds the actual class of the control plane after a successful transition.
	Current string `json:"current,omitempty"`
	// Size holds the status information about the control plane size,
	// including resource usage.
	Size *ControlPlaneSizeStatus `json:"size,omitempty"`
	// History contains a list of the control plane's class transitions and the
	// time they were observed.
	// +kubebuilder:validation:MaxItems=5
	History []*ClassTransition `json:"history,omitempty"`
}

// ControlPlaneSizeStatus defines the status of the control plane size.
type ControlPlaneSizeStatus struct {
	ResourceUsage *ResourceUsage `json:"resourceUsage,omitempty"`
}

// ResourceUsage represents the resource limits and consumption of the control plane.
type ResourceUsage struct {
	// CPU represents the CPU resource usage.
	CPU string `json:"cpu,omitempty"`
	// Memory represents the memory resource usage.
	Memory string `json:"memory,omitempty"`
	// Storage represents the storage resource usage.
	Storage string `json:"storage,omitempty"`
}

// ClassTransition represents a record of the control plane's transition to a different class.
type ClassTransition struct {
	// Class is the class to which the control plane was observed to have
	// transitioned.
	// +required
	// +kubebuilder:validation:MinLength=1
	Class string `json:"class"`
	// ObservedAt is the time at which the control plane was observed to have
	// transitioned to the associated class.
	// +required
	ObservedAt *metav1.Time `json:"observedAt"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Crossplane",type="string",JSONPath=".spec.crossplane.version"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=`.status.message`
// +kubebuilder:printcolumn:name="Class",type="string",JSONPath=".spec.class",priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories=spaces,shortName=ctp;ctps

// ControlPlane defines a managed Crossplane instance.
// +kubebuilder:validation:XValidation:rule="self.metadata.name.size() <= 63",message="control plane name cannot exceed 63 characters"
type ControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ControlPlaneSpec   `json:"spec"`
	Status ControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ControlPlaneList contains a list of ControlPlane
type ControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ControlPlane `json:"items"`
}

// GetCondition of this ControlPlane.
func (mg *ControlPlane) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return mg.Status.GetCondition(ct)
}

// GetWriteConnectionSecretToReference of this ControlPlane.
func (mg *ControlPlane) GetWriteConnectionSecretToReference() *xpv2.SecretReference {
	if mg.Spec.WriteConnectionSecretToReference == nil {
		return nil
	}
	return &xpv2.SecretReference{
		Name:      mg.Spec.WriteConnectionSecretToReference.Name,
		Namespace: mg.Spec.WriteConnectionSecretToReference.Namespace,
	}
}

// SetConditions of this ControlPlane.
func (mg *ControlPlane) SetConditions(c ...xpv2.Condition) {
	mg.Status.SetConditions(c...)
}

// SetWriteConnectionSecretToReference of this ControlPlane.
func (mg *ControlPlane) SetWriteConnectionSecretToReference(r *xpv2.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = &SecretReference{
		Name:      r.Name,
		Namespace: r.Namespace,
	}
}

// ManagedControlPlane type metadata.
var (
	// ControlPlaneKind is the kind of the ControlPlane.
	ControlPlaneKind = reflect.TypeOf(ControlPlane{}).Name()
	// ControlPlaneListKind is the kind of a list of ControlPlane.
	ControlPlaneListKind = reflect.TypeOf(ControlPlaneList{}).Name()
)

func init() {
	SchemeBuilder.Register(&ControlPlane{}, &ControlPlaneList{})
}
