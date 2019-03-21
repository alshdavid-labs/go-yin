# go-yin
Convenience utilities for http servers


```
go get -u github.com/alshdavid-sandbox/go-yin
```

```Go
package main

import (
    "fmt"
    "net/http"

    "github.com/alshdavid-sandbox/go-yin"

    "github.com/go-chi/chi"
)

func main() {
    r := chi.NewRouter()
    r.Use(yin.Logger)

    r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
        body := yin.H{}
        
        err := yin.Req(r).Body(&body)
        if err != nil {
            yin.Res(w).SendStatus(http.StatusBadRequest)
            return
        }
        
        fmt.Println(body)

        yin.Res(w).
            JSON(yin.H{
                "Hello": "world",
            })
    })

    http.ListenAndServe(":3000", r)
}
```