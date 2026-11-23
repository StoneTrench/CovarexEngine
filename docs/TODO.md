# TODO

## Platform

- [x] Document VirtualFileSystem api
- [ ] Document meta system api

## Core

- [ ] Document events SignalSource api
- [ ] Document log BracketHandler api
- [ ] Document log logging api

## Resources

- [ ] Document config api
- [ ] Document translate api
- [ ] Document registry api

## Systems

### ImGui integration

- [ ] Just like add some way to debug the ECS, the networking and the worker pool?
	- [ ] Just try ImGui <3 I've heard it's good, and I know you will forget this by tomorrow. yes, dear reader, I'am talking to myself in this document :3 Its 1am

### Physics

- [ ] Collision detection
- [ ] Rigidbody physics
- [ ] Reliable physics tick system
	- [ ] Accumulation that wont cause the system to lag behind

### ECS

- [ ] Document ECS

- [x] Raw byte component set (can be created at runtime)
	- [ ] Create component set
- [ ] Create systems from WASM at runtime
- [ ] System scheduler

### Modding

- [ ] Document module system

#### WASM abi

- [x] Parse .abi file
- [x] Translate type & function declarations into instructions
	- [x] Encode primitive
	- [x] Encode enum
	- [x] Encode slice
	- [x] Encode array
	- [x] Encode struct
	- [x] Encode pointer
- [x] Translate instructions to target language
	- [x] Guest abi implementation
	- [x] Host abi implementation

#### Dependencies

- [x] `about.json` module header file
- [x] Host module builder
- [x] Find dependency issues
	- [x] Incorrect version
	- [x] Missing dependecy
	- [x] Missing recommendency
	- [x] Missing suggestency
	- [x] Existing conflict
	- [x] Existing breakflict
	- [x] Invalid order
- [ ] Auto sort



### Networking

- [ ] Document the network protocol

- [ ] Split transmitter
	- [ ] ReliableSender
	- [ ] ReliableReceiver
- [ ] Refactor abstraction layer on top of this
	- [ ] Use signals? or channels?

### Worker pool

- [ ] Worker
	- [x] Submit jobs
	- [x] Thread safety
	- [x] Panic recovery
	- [x] Close and CloseImmediate
	- [ ] Share data across workers

