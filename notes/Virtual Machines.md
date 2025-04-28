# Virtual Machines

> [!NOTE]
> The VMs we are talking about here are different from those of VMWare or Virtualbox.

[[Between the CPU and the call stack]]

VMs == Computers built with software. VMs mimic how a computer works, and we can think of it as a _custom-built computer_

A VM can be anything (a function, a struct, an object, etc.), and what matters is _what it does_

[[Stack machine vs Register machine]]

When the VM does the **dispatching**, it means it select _an implementation for an instruction before executing it_ e.g., `PUSH`

[[Why build a VM]]

VMs are domain-specific, while computers are not. For that, our VMs only need a subset of the features a computer has to offer
