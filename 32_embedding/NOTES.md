# Notes

## Key points

- Go doesn't support inheritance. Instead, it uses composition and embedding.
- Embedding a type inside a struct promotes its fields and methods.
- Promoted fields and methods can be accessed directly from the outer type.
- Embedding can also be used with pointer types and unexported types.
- It is very common to embed interfaces too. The name of the composition is usually the composition of names (eg. the ReadWriter).

### Receiver behavior

- Promoted methods from embedded types retain their original receiver (value or pointer).
- If a pointer receiver is embedded and not initialized, calling a promoted method may panic. Always check for nil pointers when embedding pointer types.

## Takeaways

- Use embedding to reuse logic and simplify struct hierarchies.
- When embedding multiple types, be aware of field/method name clashes.
- Field and method collisions lead to non-promotion or require manual overrides.
- You can override a method in the outer struct, so that the embedded version is shadowed.
- When overriding, you can still call the inner method when needed.
- Unexported types can still be embedded and promoted from, but the inner type identifier remains inaccessible.
- Use type conversion to define alternate behaviors without modifying original types.
- Embedding a type can make the outer type satisfy the interfaces the inner type implements, but you cannot use the outer type in a function argument where the embedded type is required. 

## Try it out

- Embed a struct and access its fields and methods through the outer struct.
- Override a method and still call the embedded method from the override.
- Try embedding two types to be promoted, but both sharing the same field or method name and observe compiler errors.
- Embed a pointer to a struct and try calling its methods when uninitialized (expect a panic).

