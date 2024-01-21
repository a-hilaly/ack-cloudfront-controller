// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package function

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/cloudfront"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/cloudfront-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.CloudFront{}
	_ = &svcapitypes.Function{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.DescribeFunctionOutput
	resp, err = rm.sdkapi.DescribeFunctionWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeFunction", err)
	if err != nil {
		if reqErr, ok := ackerr.AWSRequestFailure(err); ok && reqErr.StatusCode() == 404 {
			return nil, ackerr.NotFound
		}
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "NoSuchFunctionExists" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.ETag != nil {
		ko.Status.ETag = resp.ETag
	} else {
		ko.Status.ETag = nil
	}
	if resp.FunctionSummary != nil {
		f1 := &svcapitypes.FunctionSummary{}
		if resp.FunctionSummary.FunctionConfig != nil {
			f1f0 := &svcapitypes.FunctionConfig{}
			if resp.FunctionSummary.FunctionConfig.Comment != nil {
				f1f0.Comment = resp.FunctionSummary.FunctionConfig.Comment
			}
			if resp.FunctionSummary.FunctionConfig.Runtime != nil {
				f1f0.Runtime = resp.FunctionSummary.FunctionConfig.Runtime
			}
			f1.FunctionConfig = f1f0
		}
		if resp.FunctionSummary.FunctionMetadata != nil {
			f1f1 := &svcapitypes.FunctionMetadata{}
			if resp.FunctionSummary.FunctionMetadata.CreatedTime != nil {
				f1f1.CreatedTime = &metav1.Time{*resp.FunctionSummary.FunctionMetadata.CreatedTime}
			}
			if resp.FunctionSummary.FunctionMetadata.FunctionARN != nil {
				f1f1.FunctionARN = resp.FunctionSummary.FunctionMetadata.FunctionARN
			}
			if resp.FunctionSummary.FunctionMetadata.LastModifiedTime != nil {
				f1f1.LastModifiedTime = &metav1.Time{*resp.FunctionSummary.FunctionMetadata.LastModifiedTime}
			}
			if resp.FunctionSummary.FunctionMetadata.Stage != nil {
				f1f1.Stage = resp.FunctionSummary.FunctionMetadata.Stage
			}
			f1.FunctionMetadata = f1f1
		}
		if resp.FunctionSummary.Name != nil {
			f1.Name = resp.FunctionSummary.Name
		}
		if resp.FunctionSummary.Status != nil {
			f1.Status = resp.FunctionSummary.Status
		}
		ko.Status.FunctionSummary = f1
	} else {
		ko.Status.FunctionSummary = nil
	}

	rm.setStatusDefaults(ko)
	if resp.FunctionSummary != nil {
		if resp.FunctionSummary.FunctionMetadata != nil {
			if resp.FunctionSummary.FunctionMetadata.FunctionARN != nil {
				ko.Status.ACKResourceMetadata.ARN = (*ackv1alpha1.AWSResourceName)(resp.FunctionSummary.FunctionMetadata.FunctionARN)
			}
			if resp.FunctionSummary.FunctionConfig != nil {
				if resp.FunctionSummary.FunctionConfig.Runtime != nil {
					ko.Spec.FunctionConfig.Runtime = resp.FunctionSummary.FunctionConfig.Runtime
				}
				if resp.FunctionSummary.FunctionConfig.Comment != nil {
					ko.Spec.FunctionConfig.Comment = resp.FunctionSummary.FunctionConfig.Comment
				}
			}
		}
	}
	if err := rm.setResourceAdditionalFields(ctx, ko); err != nil {
		return nil, err
	}

	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Spec.Name == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeFunctionInput, error) {
	res := &svcsdk.DescribeFunctionInput{}

	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateFunctionOutput
	_ = resp
	resp, err = rm.sdkapi.CreateFunctionWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateFunction", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.ETag != nil {
		ko.Status.ETag = resp.ETag
	} else {
		ko.Status.ETag = nil
	}
	if resp.FunctionSummary != nil {
		f1 := &svcapitypes.FunctionSummary{}
		if resp.FunctionSummary.FunctionConfig != nil {
			f1f0 := &svcapitypes.FunctionConfig{}
			if resp.FunctionSummary.FunctionConfig.Comment != nil {
				f1f0.Comment = resp.FunctionSummary.FunctionConfig.Comment
			}
			if resp.FunctionSummary.FunctionConfig.Runtime != nil {
				f1f0.Runtime = resp.FunctionSummary.FunctionConfig.Runtime
			}
			f1.FunctionConfig = f1f0
		}
		if resp.FunctionSummary.FunctionMetadata != nil {
			f1f1 := &svcapitypes.FunctionMetadata{}
			if resp.FunctionSummary.FunctionMetadata.CreatedTime != nil {
				f1f1.CreatedTime = &metav1.Time{*resp.FunctionSummary.FunctionMetadata.CreatedTime}
			}
			if resp.FunctionSummary.FunctionMetadata.FunctionARN != nil {
				f1f1.FunctionARN = resp.FunctionSummary.FunctionMetadata.FunctionARN
			}
			if resp.FunctionSummary.FunctionMetadata.LastModifiedTime != nil {
				f1f1.LastModifiedTime = &metav1.Time{*resp.FunctionSummary.FunctionMetadata.LastModifiedTime}
			}
			if resp.FunctionSummary.FunctionMetadata.Stage != nil {
				f1f1.Stage = resp.FunctionSummary.FunctionMetadata.Stage
			}
			f1.FunctionMetadata = f1f1
		}
		if resp.FunctionSummary.Name != nil {
			f1.Name = resp.FunctionSummary.Name
		}
		if resp.FunctionSummary.Status != nil {
			f1.Status = resp.FunctionSummary.Status
		}
		ko.Status.FunctionSummary = f1
	} else {
		ko.Status.FunctionSummary = nil
	}
	if resp.Location != nil {
		ko.Status.Location = resp.Location
	} else {
		ko.Status.Location = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateFunctionInput, error) {
	res := &svcsdk.CreateFunctionInput{}

	if r.ko.Spec.FunctionCode != nil {
		res.SetFunctionCode(r.ko.Spec.FunctionCode)
	}
	if r.ko.Spec.FunctionConfig != nil {
		f1 := &svcsdk.FunctionConfig{}
		if r.ko.Spec.FunctionConfig.Comment != nil {
			f1.SetComment(*r.ko.Spec.FunctionConfig.Comment)
		}
		if r.ko.Spec.FunctionConfig.Runtime != nil {
			f1.SetRuntime(*r.ko.Spec.FunctionConfig.Runtime)
		}
		res.SetFunctionConfig(f1)
	}
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newUpdateRequestPayload(ctx, desired, delta)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.UpdateFunctionOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateFunctionWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateFunction", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.ETag != nil {
		ko.Status.ETag = resp.ETag
	} else {
		ko.Status.ETag = nil
	}
	if resp.FunctionSummary != nil {
		f1 := &svcapitypes.FunctionSummary{}
		if resp.FunctionSummary.FunctionConfig != nil {
			f1f0 := &svcapitypes.FunctionConfig{}
			if resp.FunctionSummary.FunctionConfig.Comment != nil {
				f1f0.Comment = resp.FunctionSummary.FunctionConfig.Comment
			}
			if resp.FunctionSummary.FunctionConfig.Runtime != nil {
				f1f0.Runtime = resp.FunctionSummary.FunctionConfig.Runtime
			}
			f1.FunctionConfig = f1f0
		}
		if resp.FunctionSummary.FunctionMetadata != nil {
			f1f1 := &svcapitypes.FunctionMetadata{}
			if resp.FunctionSummary.FunctionMetadata.CreatedTime != nil {
				f1f1.CreatedTime = &metav1.Time{*resp.FunctionSummary.FunctionMetadata.CreatedTime}
			}
			if resp.FunctionSummary.FunctionMetadata.FunctionARN != nil {
				f1f1.FunctionARN = resp.FunctionSummary.FunctionMetadata.FunctionARN
			}
			if resp.FunctionSummary.FunctionMetadata.LastModifiedTime != nil {
				f1f1.LastModifiedTime = &metav1.Time{*resp.FunctionSummary.FunctionMetadata.LastModifiedTime}
			}
			if resp.FunctionSummary.FunctionMetadata.Stage != nil {
				f1f1.Stage = resp.FunctionSummary.FunctionMetadata.Stage
			}
			f1.FunctionMetadata = f1f1
		}
		if resp.FunctionSummary.Name != nil {
			f1.Name = resp.FunctionSummary.Name
		}
		if resp.FunctionSummary.Status != nil {
			f1.Status = resp.FunctionSummary.Status
		}
		ko.Status.FunctionSummary = f1
	} else {
		ko.Status.FunctionSummary = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
	delta *ackcompare.Delta,
) (*svcsdk.UpdateFunctionInput, error) {
	res := &svcsdk.UpdateFunctionInput{}

	if r.ko.Spec.FunctionCode != nil {
		res.SetFunctionCode(r.ko.Spec.FunctionCode)
	}
	if r.ko.Spec.FunctionConfig != nil {
		f1 := &svcsdk.FunctionConfig{}
		if r.ko.Spec.FunctionConfig.Comment != nil {
			f1.SetComment(*r.ko.Spec.FunctionConfig.Comment)
		}
		if r.ko.Spec.FunctionConfig.Runtime != nil {
			f1.SetRuntime(*r.ko.Spec.FunctionConfig.Runtime)
		}
		res.SetFunctionConfig(f1)
	}
	if r.ko.Status.ETag != nil {
		res.SetIfMatch(*r.ko.Status.ETag)
	}
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteFunctionOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteFunctionWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteFunction", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteFunctionInput, error) {
	res := &svcsdk.DeleteFunctionInput{}

	if r.ko.Status.ETag != nil {
		res.SetIfMatch(*r.ko.Status.ETag)
	}
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Function,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	// No terminal_errors specified for this resource in generator config
	return false
}