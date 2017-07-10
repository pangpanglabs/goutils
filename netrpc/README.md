# netrpc

Reconnectable net/rpc client.

## Getting Started

```
func CalculateCart(addr string, cart *Cart) (*Cart, error) {
client := netprc.NewClient(addr)
var result Cart
if err := client.Call("offer.CalculateCart", cart, &result); err != nil {
        return nil, err
}
return &result, nil
```