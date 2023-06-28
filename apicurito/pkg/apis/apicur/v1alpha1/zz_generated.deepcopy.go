// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Apicurito) DeepCopyInto(out *Apicurito) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Apicurito.
func (in *Apicurito) DeepCopy() *Apicurito {
	if in == nil {
		return nil
	}
	out := new(Apicurito)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Apicurito) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApicuritoList) DeepCopyInto(out *ApicuritoList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Apicurito, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApicuritoList.
func (in *ApicuritoList) DeepCopy() *ApicuritoList {
	if in == nil {
		return nil
	}
	out := new(ApicuritoList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ApicuritoList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApicuritoSpec) DeepCopyInto(out *ApicuritoSpec) {
	*out = *in
	if in.ResourcesUI != nil {
		in, out := &in.ResourcesUI, &out.ResourcesUI
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.ResourcesGenerator != nil {
		in, out := &in.ResourcesGenerator, &out.ResourcesGenerator
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApicuritoSpec.
func (in *ApicuritoSpec) DeepCopy() *ApicuritoSpec {
	if in == nil {
		return nil
	}
	out := new(ApicuritoSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApicuritoStatus) DeepCopyInto(out *ApicuritoStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApicuritoStatus.
func (in *ApicuritoStatus) DeepCopy() *ApicuritoStatus {
	if in == nil {
		return nil
	}
	out := new(ApicuritoStatus)
	in.DeepCopyInto(out)
	return out
}
