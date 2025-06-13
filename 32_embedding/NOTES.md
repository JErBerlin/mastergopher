# Notes

## Key points

- Go doesn't support inheritance. Instead, it uses composition and embedding.
- Embedding a type inside a struct promotes its fields and methods.
- Promoted fields and methods can be accessed directly from the outer type.
- Field and method collisions lead to non-promotion or require manual overrides.
- Method overrides can still access the original method from the embedded type.
- Embedding can also be used with pointer types and unexported types.
- When embedding pointers, nil pointer dereferences can cause runtime panics.

### Receiver behavior

- Promoted methods from embedded types retain their original receiver (value or pointer).
- If a pointer receiver is embedded and not initialized, calling a promoted method may panic.
- If you override a method in the outer struct, the embedded version is shadowed.

## Takeaways

- Use embedding to reuse logic and simplify struct hierarchies.
- When embedding multiple types, avoid field/method name clashes.
- When overriding, you can still call the inner method when needed.
- Unexported types can still be embedded and promoted from, but the inner type identifier remains inaccessible.
- Use type conversion to define alternate behaviors without modifying original types.
- Always check for nil pointers when embedding pointer types and calling a (possibly) promoted method.
- Embedding a type does not make the outer type satisfy the interfaces it is implementing, even if the right methods are promoted. 

## Try it out

- Embed a struct and access its fields and methods through the outer struct.
- Override a method and still call the embedded method from the override.
- Try embedding two types to be promoted, but both sharing the same field or method name and observe compiler errors.
- Embed a pointer to a struct and try calling its methods when uninitialized (expect a panic).

