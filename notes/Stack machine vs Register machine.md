# Stack machine vs Register machine

One major design decision: Whether the machine is a _stack machine_ to do its computation or a _register machine_ (virtual one!)

While the stack machine is simpler to build, there is a performance limit on the stack since we need to push/pop a lot of things on/off the stack

The register machine, while harder to build, takes advantage of the registers and thus much denser compared to the stack machine. Instructions can refer to the registers directly
