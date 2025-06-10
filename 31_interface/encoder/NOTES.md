# Notes

## Key points

- After assigning to an interface, this holds a pair: the dynamic type and the dynamic value. If neither is set, the interface is `nil`. Using an uninitialized interface will not fail at compile time but can panic at runtime. 
- Interfaces become very useful when backed by a concrete type that implements the required methods.
- Receiver types matter: pointer receivers and value receivers interact differently with interfaces. Value receivers are more flexible in how they satisfy interfaces, while pointer receivers are required when methods need to mutate the receiver.

## Receiver behavior

- A method with a **pointer receiver** (`func (t *T) Method()`) can only be called on a pointer.
  - Structs with such methods will only satisfy an interface if passed by reference.
- A method with a **value receiver** (`func (t T) Method()`) is more flexible and allows both value and pointer usage.
  - This means a value type can satisfy the interface directly, even if passed by pointer.

## Type assertions and type switches

- A `switch` on the result of a type assertion (`switch v := val.(type)`) allows branching based on the dynamic type stored in an interface.
- This is commonly used in generic functions (e.g., serializers or encoders) to delegate logic based on the input's actual type.
- You can switch both on concrete types (e.g. `string`, `int`, `map[string]string`) and on interfaces (e.g. `json.Marshaler`).

## Takeaways

- Always ensure embedded interfaces or interface fields are initialized before use.
- When building APIs that accept general values (`interface{}`), prefer using type switches to direct specific logic per type.
- Prefer pointer receivers when:                                                                                   - You need to mutate the struct.
  - You want to avoid copying large structs.
- Prefer value receivers when:                                                                                     - The method does not modify the struct.
- Use `json.Marshaler` or similar interfaces with type switches to let types define custom encoding logic.

Example of a type switch:

(Go)  
`switch v := input.(type) {`  
`case CustomEncoder: ...`  
`case string: ...`  
`case int: ...`  
`default: ...`  
`}`

This pattern is widespread in libraries and generic helpers.

## Try it out

The `encode` function is meant to take a general value and produce a formatted output depending on its type or capabilities.

- Try passing a type not yet handled (e.g. `float64` or a custom struct) to observe the fallback behavior.
- Create your own type and implement a custom encoder interface (e.g. with `MarshalJSON`) to test interface-based dispatch.

