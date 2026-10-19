# Design overview

## Types

### Primitive

+ Just a number.
+ Gets sent directly and maybe cast to a different number along the way.

### Plex

+ Product type, like an array or struct.
+ Plexes can be large or small (depending on the number of elements if flattened, if it's more than 8 elements, then it is large).

### Pointer

+ Pointers can be constant or mutable.

## Ownership & mutability

A value passed by a pointer has:
+ Mutability (more restrictive state is propagated):
	+ Constant: read only (propagated).
	+ Mutable: read and write (only this pointer).
+ Ownership (use-site modifier, can only appear for function parameters):
	+ Given: recipient can free, sender no longer owns it.
	+ Borrowed: recipient uses it during call, sender can free.

+ Function parameters are `borrowed` by the `host` from the `guest`.
+ Function returns are `given` to the `guest` by the `host`

## Large plex

+ A `plex` is considered large, if the number of total fields (including the children's fields recursively), excede the `threshold` (TODO: figure out a threshold), then instead of getting flattened, the `plex` is passed by `pointer`.
+ If a `plex` is passed by `pointer` (because it's large), the `pointer` will always be constant.
+ Flattening means, to unwrap the fields/members of a `plex` and treat them as individual `primitives`.

## Pointers
+ If a `pointer` is constant it will always get copied at the `host`/`guest` boundary, including the nested `pointers` recursively. If the parent `pointer` is constant all it's child `pointers` will also be treated as constant.
+ Copying is always done by the `host`, never the `guest`, it is also syncronous.
+ Both the `host` and `guest` can allocate inside the `guest`.
+ Only the `guest` can free its own memory (this means that only the `guest` can free, the `host` can not free).
+ If the `pointer` is not constant, that means the `host` or `guest` can write to it.
+ If the `host` can write to it, the `guest` must free based on the documentation of the function.

## Guest writing to Host (`G->H`)

### Primitive

+ A `primitive` is passed directly.
+ If a `pointer` is passed, the `host` copies the data from the memory of the `guest`. Copying happens syncronously if the `pointer` is constant, this way the `host`/`guest` cannot modify it by accident.

### Plex

+ If it's small, the `plex` is flattened.
+ If it's large, a `pointer` to the `plex` is passed to the `host`. At which point it can be treated as a `primitive`.

## Host writing to Guest (`H->G`)

### Primitive

+ A `primitive` is passed directly.
+ If a `pointer` is passed, the `host` allocates memory inside the `guest`, and copies the data into the allocated memory, then passes the `pointer` to the `guest`.

### Plex

+ If it's small, the `plex` is flattened.
+ If it's large, the `host` allocates memory inside the `guest`, and copies the `plex` into the allocated memory, then passes the `pointer` to `guest`.


