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