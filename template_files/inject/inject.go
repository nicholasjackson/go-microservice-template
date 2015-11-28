/*
 * Copyright 2013 Google Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package inject

import (
	"fmt"
	"reflect"
)

// Context that is passed to scopes.
type Context interface{}

type Singleton struct{}

// Key used to uniquely identify a binding.
type Key interface{}
type Tag interface{}

// Type used to identify a tagged type binding.
type TaggedKey struct {
	Key
	Tag
}

/*
	Signature for provider functions. Provider functions are used to dynamically allocate an instance
	at run-time.
*/
type Provider func(Context, Container) interface{}

/*
	Injector aggregates binding configuration and creates Containers based on that configuration.
	Binding configuration is defined by a Key that consists of a type and an optional tag.

	A child injector may be used to create bindings that are intended to be used only by part
	of the system. When a Container is created from a child injector, the bindings of the
	parent Injector are also available.

	Keys may only be bound once across an Injector and all of its descendant Injectors (children
	and their children). The Injector will panic if an attempt is made to rebind an already-bound
	Key.

	In order to look up a bound type, use the CreateContainer() method and call the appropriate
	methods on the returned Container.
*/
type Injector interface {
	// Binds a type to a Provider function.
	Bind(Key, Provider)

	// Binds a type to a single instance.
	BindInstance(Key, interface{})

	// Binds a key to a Provider function, caching it within the specified scope.
	BindInScope(Key, Provider, Tag)

	// Binds a key to a single instance.
	BindInstanceInScope(Key, interface{}, Tag)

	// Binds a scope to a tag.
	BindScope(Scope, Tag)

	// Binds a tagged type to a Provider function.
	BindTagged(Key, Tag, Provider)

	// Binds a tagged type to a Provider function.
	BindTaggedInScope(Key, Tag, Provider, Tag)

	// Binds a tagged type to a single instance.
	BindTaggedInstance(Key, Tag, interface{})

	// Binds a tagged type to a single instance.
	BindTaggedInstanceInScope(Key, Tag, interface{}, Tag)

	// Creates a child injector that can bind additional types not available from this Injector.
	CreateChildInjector() Injector

	// Creates a Container that can be used to retrieve instance objects from the Injector.
	CreateContainer() Container

	// Exposes a type to its parent injector.
	Expose(Key)

	// Exposes a type to its parent injector.
	ExposeAndRename(Key, Key)

	// Exposes a type to its parent injector.
	ExposeTaggedAndRename(Key, Tag, Key)

	// Exposes a type to its parent injector.
	ExposeAndRenameTagged(Key, Key, Tag)

	// Exposes a type to its parent injector.
	ExposeAndTag(Key, Tag)

	// Exposes a type to its parent injector.
	ExposeTaggedAndRenameTagged(Key, Tag, Key, Tag)

	// Exposes a tagged type to its parent injector.
	ExposeTagged(Key, Tag)

	// Wraps a Provider to cache in a given scope.
	Scope(key Key, provider Provider, scopeTag Tag) Provider

	// Gets the binding for a key, searching the current injector and all ancestor injectors.
	getBinding(Key) (binding, bool)

	/*
		Searches the parent injector for the key, continuing to search upward until the
		root injector is found.
	*/
	findAncestorBinding(Key) (binding, bool)
}

// The context holds all the keys used by a given object.
type keyset map[Key]Key

type binding struct {
	injector *injector
	provider Provider
}

// Bindings for each key in the injector.
type bindings map[Key]binding

type scopes map[Tag]Scope

type injector struct {
	// The bindings present in this injector.
	bindings

	// Registered scopes (shared among all injectors)
	scopes

	// The parent injector. See getBinding(), findAncestorBinding().
	parent *injector

	// A pointer to the keyset for this injector and all ancestor and descendant injectors.
	keyset
}

func CreateInjector() Injector {
	singleton := singletonscope{make(map[Key]interface{})}
	scopes := make(scopes)
	scopes[Singleton{}] = &singleton
	return &injector{
		bindings: make(map[Key]binding),
		scopes:   scopes,
		parent:   nil,
		keyset:   make(keyset),
	}
}

// Creates a child injector that can contain bindings not available to the parent injector.
func (this *injector) CreateChildInjector() Injector {
	child := injector{
		bindings: make(map[Key]binding),
		scopes:   this.scopes,
		parent:   this,
		keyset:   this.keyset,
	}

	return &child
}

func (this injector) Bind(key Key, provider Provider) {
	if _, exists := this.bindings[key]; exists {
		panic(fmt.Sprintf("%s is already bound.", key))
	}

	if _, exists := this.findAncestorBinding(key); exists {
		panic(fmt.Sprintf("%s is already bound in an ancestor injector.", key))
	}

	this.keyset[key] = key
	this.bindings[key] = binding{&this, provider}
}

func (this injector) Scope(key Key, provider Provider, scopeTag Tag) Provider {
	var scopes = this.scopes
	if scope, exists := scopes[scopeTag]; exists {
		return scope.Scope(key, provider)
	}
	panic(fmt.Sprintf("Scope tag '%s' is not bound", scopeTag))
}

func (this injector) BindInScope(key Key, provider Provider, scopeTag Tag) {
	this.Bind(key, this.Scope(key, provider, scopeTag))
}

func (this injector) BindInstance(key Key, instance interface{}) {
	this.Bind(key, func(context Context, container Container) interface{} { return instance })
}

func (this injector) BindInstanceInScope(key Key, value interface{}, scopeTag Tag) {
	this.BindInScope(key, func(context Context, container Container) interface{} { return value }, scopeTag)
}

func (this injector) BindTaggedInScope(bindingType Key, tag Tag, provider Provider, scopeTag Tag) {
	this.BindInScope(TaggedKey{bindingType, tag}, provider, scopeTag)
}

func (this injector) BindTagged(instanceType Key, tag Tag, provider Provider) {
	this.Bind(TaggedKey{instanceType, tag}, provider)
}

func (this injector) BindTaggedInstance(instanceType Key, tag Tag,
	instance interface{}) {
	this.BindInstance(TaggedKey{instanceType, tag}, instance)
}

func (this injector) BindTaggedInstanceInScope(bindingType Key, tag Tag, value interface{}, scopeTag Tag) {
	this.BindInstanceInScope(TaggedKey{bindingType, tag}, value, scopeTag)
}

// Creates a Container that is used to request values during object creation.
func (this injector) CreateContainer() Container {
	return container{
		&this,
		make(keyset),
		nil,
		make(map[*injector]*container),
	}
}

func (this injector) Expose(key Key) {
	this.ExposeAndRename(key, key)
}

func (this injector) ExposeTagged(key Key, tag Tag) {
	this.ExposeAndRename(TaggedKey{key, tag}, TaggedKey{key, tag})
}

func (this injector) ExposeAndTag(key Key, tag Tag) {
	this.ExposeAndRename(key, TaggedKey{key, tag})
}

func (this injector) ExposeTaggedAndRename(key Key, tag Tag, parentKey Key) {
	this.ExposeAndRename(TaggedKey{key, tag}, parentKey)
}

func (this injector) ExposeTaggedAndRenameTagged(key Key, tag Tag, parentKey Key, parentTag Tag) {
	this.ExposeAndRename(TaggedKey{key, tag}, TaggedKey{parentKey, parentTag})
}

func (this injector) ExposeAndRenameTagged(key Key, parentKey Key, parentTag Tag) {
	this.ExposeAndRename(key, TaggedKey{parentKey, parentTag})
}

func (this injector) ExposeAndRename(childKey Key, parentKey Key) {
	if this.parent == nil {
		panic(fmt.Sprintf("No parent injector available when exposing %s.", childKey))
	}
	if _, exists := this.bindings[childKey]; !exists {
		panic(fmt.Sprintf("No binding for %s is present in the child injector.", childKey))
	}
	if _, exists := this.findAncestorBinding(parentKey); exists {
		panic(fmt.Sprintf("A binding for %s already exists. It could come from another child injector or an ancestor injector.", parentKey))
	}

	this.parent.bindings[parentKey] = this.bindings[childKey]
}

func (this injector) getBinding(key Key) (binding, bool) {
	binding, ok := this.bindings[key]
	return binding, ok
}

func (this injector) findAncestorBinding(key Key) (binding, bool) {
	parent := this.parent
	for parent != nil {
		if binding, ok := this.parent.getBinding(key); ok {
			return binding, ok
		}

		parent = parent.parent
	}

	return binding{}, false
}

/*
	Container provides access to the bindings configured in an Injector. All bindings are available
	as a Provider or as a value. A new Container should be used for each injected type. A Container
	will panic if a key is looked up more than once. This behavior is intended to detect and prevent
	cycles in depedencies.

	For example, suppose you have a type A that gets an instance of type B that in turn relies
	on type A again (A -> B -> A). The types would have a structure like this:

	type A struct {
		B
	}

	type B struct {
		A
	}

	func ConfigureInjector(injector inject.Injector) {
		injector.Bind(reflect.TypeOf(A(nil)), func (container Container) interface{} {
			return A { createB(container) }
		}
		injector.Bind(reflect.TypeOf(B(nil)), func (container Container) interface{} {
			return B { createA(container) }
		}
	}

	func createA(container inject.Container) {
		return A { B: container.GetInstanceForKey(reflect.TypeOf(A(nil))) }
	}

	func createB(container goos.Container) {
		return B { A: container.GetInstanceForKey(reflect.TypeOf(B(nil)) }
	}
*/
type Container interface {
	// Returns an instance of the type.
	GetInstance(Context, Key) interface{}

	// Returns an instance of the type tagged with the tag.
	GetTaggedInstance(Context, Key, Tag) interface{}

	// Returns a Provider that can return an instance of the type.
	GetProvider(Key) Provider

	// Returns a Provider that can return an instance of the type tagged with the tag.
	GetTaggedProvider(Key, Tag) Provider
}

type container struct {
	// The injector holding the bindings available to the container.
	injector *injector

	// The invocation keyset, holding all the previous requests to prevent duplicate requests.
	keyset

	// The parent container; used for exposed child bindings.
	parent *container

	// Child container for each injector
	children map[*injector]*container
}

func createChildProvider(parent *container, binding binding) Provider {
	var childContainer *container

	if bindingContainer, ok := parent.children[binding.injector]; ok {
		childContainer = bindingContainer
	} else {
		container := container{
			binding.injector,
			make(keyset),
			parent,
			parent.children,
		}
		childContainer = &container
		parent.children[binding.injector] = childContainer
	}

	return func(context Context, container Container) interface{} {
		return binding.provider(context, childContainer)
	}
}

// Returns a Provider that can create an instance of the type bound to the key.
func (this container) GetProvider(key Key) Provider {
	if _, exists := this.keyset[key]; exists {
		panic(fmt.Sprintf("Already looked up %s (%+v). Is there a cycle of dependencies?", key, reflect.TypeOf(key)))
	}

	this.keyset[key] = key

	if binding, ok := this.injector.bindings[key]; ok {
		if binding.injector == this.injector {
			return binding.provider
		} else {
			return createChildProvider(&this, binding)
		}
	}

	if binding, ok := this.injector.findAncestorBinding(key); ok {
		return createChildProvider(&this, binding)
	}

	panic(fmt.Sprintf("Unable to find %s in injector", key))
}

// Returns a Provider that can create an instance of the instanceType tagged with tag.
func (this container) GetTaggedProvider(instanceType Key, tag Tag) Provider {
	return this.GetProvider(TaggedKey{instanceType, tag})
}

// Returns an instance of the type bound to the key.
func (this container) GetInstance(context Context, key Key) interface{} {
	return this.GetProvider(key)(context, this)
}

// Returns an instance of the instanceType tagged with tag.
func (this container) GetTaggedInstance(context Context, instanceType Key, tag Tag) interface{} {
	return this.GetInstance(context, TaggedKey{instanceType, tag})
}

func (this TaggedKey) String() string {
	if this.Tag == nil {
		return fmt.Sprintf("%v<%s>", reflect.TypeOf(this.Key), reflect.TypeOf(this.Tag))
	}
	return fmt.Sprintf("%v<%s(%v)>", reflect.TypeOf(this.Key), reflect.TypeOf(this.Tag), this.Tag)
}

type simplescope struct {
	name   string
	values map[Context]map[Key]interface{}
}

type Scope interface {
	Scope(Key, Provider) Provider
}

type scopeKey struct {
	Context
	Key
}

type SimpleScope interface {
	Scope(Key, Provider) Provider
	Enter(Context)
	Exit(Context)
}

func (this *simplescope) Enter(context Context) {
	this.values[context] = make(map[Key]interface{})
}

func (this *simplescope) Exit(context Context) {
	if _, exists := this.values[context]; exists {
		delete(this.values, context)
	} else {
		panic(fmt.Sprintf("Already out of context when existing scope %v", this))
	}
}

func (this *simplescope) Scope(key Key, provider Provider) Provider {
	return func(context Context, container Container) interface{} {
		if scope, exists := this.values[context]; exists {
			if value, exists := scope[key]; exists {
				return value
			}

			value := provider(context, container)
			scope[key] = value

			return value
		}
		panic(fmt.Sprintf("Attempt to access %s outside of scope %s. %d scopes are active.", key, this.name, len(this.values)))
	}
}

func CreateSimpleScope() SimpleScope {
	scope := simplescope{name: "SimpleScope", values: make(map[Context]map[Key]interface{})}
	return &scope
}

func CreateSimpleScopeWithName(name string) SimpleScope {
	scope := simplescope{name: name, values: make(map[Context]map[Key]interface{})}
	return &scope
}

func (this injector) BindScope(scope Scope, scopeTag Tag) {
	if _, exists := this.scopes[scopeTag]; exists {
		panic(fmt.Sprintf("Scope is already bound for tag '%s'", scopeTag))
	}
	this.scopes[scopeTag] = scope
}

type singletonscope struct {
	values map[Key]interface{}
}

func (this *singletonscope) Enter(context Context) {
	panic("You're always in singletonscope. Do not try to enter this scope.")
}

func (this *singletonscope) Exit(context Context) {
	panic("You're always in singletonscope. Do not try to exit this scope.")
}

func (this *singletonscope) Scope(key Key, provider Provider) Provider {
	return func(context Context, container Container) interface{} {
		if value, exists := this.values[key]; exists {
			return value
		}

		value := provider(context, container)
		this.values[key] = value

		return value
	}
}
