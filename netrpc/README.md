# netrpc

Reconnectable net/rpc client.

## Getting Started

```golang
var client = netprc.NewClient(addr)
func CalculateCart(cart *Cart) (*Cart, error) {
        var result Cart
        if err := client.Call("offer.CalculateCart", cart, &result); err != nil {
                return nil, err
        }
        return &result, nil
}
```