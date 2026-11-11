# TODO

## WASM abstraction

- [ ] Write an abstraction layer ontop of WASM modules.
	- [ ] Create module.
		- [ ] From bytes.
		- [ ] From file.
	- [ ] Assign module config.
	- [ ] Allow shared memory between modules.
	- [ ] Module instantiation.
		- [ ] Requires HostApi object.
	- [ ] Link functions between modules (without shared memory).
	- [ ] If memory is shared, does malloc also get shared?
	<!-- cant move binding utils cause they cannot be called then unfortunately -->
	<!-- - [ ] Move host binding utility functions here
			- [ ] Panic
			- [ ] Panicf
			- [ ] Panics
		- [ ] Read
			- [ ] ReadString
			- [ ] ReadPlex
			- [ ] ReadSlice
			- [ ] ReadPointer
			- [ ] ReadAnyFixed
		- [ ] Write
			- [ ] WriteString
			- [ ] WritePlex
			- [ ] WriteSlice
			- [ ] WritePointer
			- [ ] WriteAnyFixed
	- [ ] Binding generator should  -->

- [ ] HostApi object
	- [ ] List of HostBindings to apply to a WASM module instance.

- [x] HostBinding
```go
type HostBinding struct {
	Name string
	Fn   interface{}
}
```