// Copyright 2019 Istio Authors
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

package collection_test

import (
	"testing"

	_ "github.com/gogo/protobuf/types"
	. "github.com/onsi/gomega"

	"istio.io/istio/galley/pkg/config/schema/collection"
	"istio.io/istio/galley/pkg/config/schema/resource"
)

var (
	emptyResource = resource.Builder{
		Kind:         "Empty",
		Plural:       "empties",
		ProtoPackage: "github.com/gogo/protobuf/types",
		Proto:        "google.protobuf.Empty",
	}.MustBuild()

	structResource = resource.Builder{
		Kind:         "Struct",
		Plural:       "structs",
		ProtoPackage: "github.com/gogo/protobuf/types",
		Proto:        "google.protobuf.Struct",
	}.MustBuild()
)

func TestSchemas_Basic(t *testing.T) {
	g := NewGomegaWithT(t)

	s := collection.Builder{
		Name:     "foo",
		Resource: emptyResource,
	}.MustBuild()

	schemas := collection.SchemasFor(s)
	g.Expect(schemas.All()).To(HaveLen(1))
	g.Expect(schemas.All()[0]).To(Equal(s))
}

func TestSchemas_MustAdd(t *testing.T) {
	g := NewGomegaWithT(t)
	defer func() {
		r := recover()
		g.Expect(r).To(BeNil())
	}()
	b := collection.NewSchemasBuilder()

	s := collection.Builder{
		Name:     "foo",
		Resource: emptyResource,
	}.MustBuild()
	b.MustAdd(s)
}

func TestSchemas_MustRegister_Panic(t *testing.T) {
	g := NewGomegaWithT(t)
	defer func() {
		r := recover()
		g.Expect(r).NotTo(BeNil())
	}()
	b := collection.NewSchemasBuilder()

	s := collection.Builder{
		Name:     "foo",
		Resource: emptyResource,
	}.MustBuild()
	b.MustAdd(s)
	b.MustAdd(s)
}

func TestSchemas_Find(t *testing.T) {
	g := NewGomegaWithT(t)

	s := collection.Builder{
		Name:     "foo",
		Resource: emptyResource,
	}.MustBuild()

	schemas := collection.SchemasFor(s)

	s2, found := schemas.Find("foo")
	g.Expect(found).To(BeTrue())
	g.Expect(s2).To(Equal(s))

	_, found = schemas.Find("zoo")
	g.Expect(found).To(BeFalse())
}

func TestSchemas_MustFind(t *testing.T) {
	g := NewGomegaWithT(t)
	defer func() {
		r := recover()
		g.Expect(r).To(BeNil())
	}()

	b := collection.NewSchemasBuilder()

	s := collection.Builder{
		Name:     "foo",
		Resource: emptyResource,
	}.MustBuild()

	b.MustAdd(s)
	schemas := b.Build()

	s2 := schemas.MustFind("foo")
	g.Expect(s2).To(Equal(s))
}

func TestSchemas_MustFind_Panic(t *testing.T) {
	g := NewGomegaWithT(t)
	defer func() {
		r := recover()
		g.Expect(r).NotTo(BeNil())
	}()

	b := collection.NewSchemasBuilder()

	s := collection.Builder{
		Name:     "foo",
		Resource: emptyResource,
	}.MustBuild()

	b.MustAdd(s)
	schemas := b.Build()

	_ = schemas.MustFind("zoo")
}

func TestSchema_FindByGroupAndKind(t *testing.T) {
	g := NewGomegaWithT(t)

	s := collection.Builder{
		Name: "foo",
		Resource: resource.Builder{
			ProtoPackage: "github.com/gogo/protobuf/types",
			Proto:        "google.protobuf.Empty",
			Group:        "mygroup",
			Kind:         "Empty",
			Plural:       "empties",
		}.MustBuild(),
	}.MustBuild()

	schemas := collection.SchemasFor(s)

	s2, found := schemas.FindByGroupAndKind("mygroup", "Empty")
	g.Expect(found).To(BeTrue())
	g.Expect(s2).To(Equal(s))

	_, found = schemas.FindByGroupAndKind("bad", "bad")
	g.Expect(found).To(BeFalse())
}

func TestSchema_MustFindByGroupAndKind(t *testing.T) {
	g := NewGomegaWithT(t)
	b := collection.NewSchemasBuilder()

	s := collection.Builder{
		Name: "foo",
		Resource: resource.Builder{
			ProtoPackage: "github.com/gogo/protobuf/types",
			Proto:        "google.protobuf.Empty",
			Group:        "mygroup",
			Kind:         "Empty",
			Plural:       "empties",
		}.MustBuild(),
	}.MustBuild()

	b.MustAdd(s)
	schemas := b.Build()

	s2 := schemas.MustFindByGroupAndKind("mygroup", "Empty")
	g.Expect(s2).To(Equal(s))
}

func TestSchema_MustFindByGroupAndKind_Panic(t *testing.T) {
	g := NewGomegaWithT(t)

	defer func() {
		r := recover()
		g.Expect(r).NotTo(BeNil())
	}()

	schemas := collection.NewSchemasBuilder().Build()
	_ = schemas.MustFindByGroupAndKind("mygroup", "Empty")
}

func TestSchema_MustFindByKind(t *testing.T) {
	g := NewGomegaWithT(t)
	b := collection.NewSchemasBuilder()

	s := collection.Builder{
		Name: "foo",
		Resource: resource.Builder{
			ProtoPackage: "github.com/gogo/protobuf/types",
			Proto:        "google.protobuf.Empty",
			Group:        "mygroup",
			Kind:         "Empty",
			Plural:       "empties",
		}.MustBuild(),
	}.MustBuild()

	b.MustAdd(s)
	schemas := b.Build()

	s2 := schemas.MustFindByKind("Empty")
	g.Expect(s2).To(Equal(s))
}

func TestSchema_MustFindByKind_Panic(t *testing.T) {
	g := NewGomegaWithT(t)

	defer func() {
		r := recover()
		g.Expect(r).NotTo(BeNil())
	}()

	schemas := collection.NewSchemasBuilder().Build()
	_ = schemas.MustFindByKind("Empty")
}

func TestSchemas_CollectionNames(t *testing.T) {
	g := NewGomegaWithT(t)
	b := collection.NewSchemasBuilder()

	s1 := collection.Builder{
		Name:     "foo",
		Resource: emptyResource,
	}.MustBuild()
	s2 := collection.Builder{
		Name:     "bar",
		Resource: emptyResource,
	}.MustBuild()
	b.MustAdd(s1)
	b.MustAdd(s2)

	names := b.Build().CollectionNames()
	expected := collection.Names{collection.NewName("bar"), collection.NewName("foo")}
	g.Expect(names).To(Equal(expected))
}

func TestSchemas_Kinds(t *testing.T) {
	g := NewGomegaWithT(t)

	s := collection.SchemasFor(
		collection.Builder{
			Name:     "foo",
			Resource: emptyResource,
		}.MustBuild(),
		collection.Builder{
			Name:     "bar",
			Resource: emptyResource,
		}.MustBuild(),
		collection.Builder{
			Name:     "baz",
			Resource: structResource,
		}.MustBuild())

	actual := s.Kinds()
	expected := []string{emptyResource.Kind(), structResource.Kind()}
	g.Expect(actual).To(Equal(expected))
}

func TestSchemas_Validate(t *testing.T) {
	cases := []struct {
		name        string
		schemas     []collection.Schema
		expectError bool
	}{
		{
			name: "valid",
			schemas: []collection.Schema{
				collection.Builder{
					Name: "foo",
					Resource: resource.Builder{
						Kind:   "Empty1",
						Plural: "Empty1s",
						Proto:  "google.protobuf.Empty",
					}.MustBuild(),
				}.MustBuild(),
				collection.Builder{
					Name: "bar",
					Resource: resource.Builder{
						Kind:   "Empty2",
						Plural: "Empty2s",
						Proto:  "google.protobuf.Empty",
					}.MustBuild(),
				}.MustBuild(),
			},
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			b := collection.NewSchemasBuilder()
			for _, s := range c.schemas {
				b.MustAdd(s)
			}
			err := b.Build().Validate()
			if c.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestSchemas_Validate_Error(t *testing.T) {
	g := NewGomegaWithT(t)
	b := collection.NewSchemasBuilder()

	s1 := collection.Builder{
		Name: "foo",
		Resource: resource.Builder{
			Kind:         "Zoo",
			ProtoPackage: "github.com/gogo/protobuf/types",
			Proto:        "zoo",
		}.BuildNoValidate(),
	}.MustBuild()
	b.MustAdd(s1)

	err := b.Build().Validate()
	g.Expect(err).NotTo(BeNil())
}

func TestSchemas_ForEach(t *testing.T) {
	schemas := collection.SchemasFor(
		collection.Builder{
			Name:     "foo",
			Resource: emptyResource,
		}.MustBuild(),
		collection.Builder{
			Name:     "bar",
			Resource: emptyResource,
		}.MustBuild(),
	)

	cases := []struct {
		name     string
		expected []string
		actual   func() []string
	}{
		{
			name:     "all",
			expected: []string{"foo", "bar"},
			actual: func() []string {
				a := make([]string, 0)
				schemas.ForEach(func(s collection.Schema) bool {
					a = append(a, s.Name().String())
					return false
				})
				return a
			},
		},
		{
			name:     "exit early",
			expected: []string{"foo"},
			actual: func() []string {
				a := make([]string, 0)
				schemas.ForEach(func(s collection.Schema) bool {
					a = append(a, s.Name().String())
					return true
				})
				return a
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			actual := c.actual()
			g.Expect(actual).To(Equal(c.expected))
		})
	}
}

func TestSchemas_Remove(t *testing.T) {
	g := NewGomegaWithT(t)

	foo := collection.Builder{
		Name:     "foo",
		Resource: emptyResource,
	}.MustBuild()
	bar := collection.Builder{
		Name:     "bar",
		Resource: emptyResource,
	}.MustBuild()
	baz := collection.Builder{
		Name:     "baz",
		Resource: emptyResource,
	}.MustBuild()

	schemas := collection.SchemasFor(foo, bar)
	g.Expect(schemas.Remove(bar.Name())).To(Equal(collection.SchemasFor(foo)))
	g.Expect(schemas.Remove(foo.Name(), bar.Name(), baz.Name())).To(Equal(collection.SchemasFor()))
}
