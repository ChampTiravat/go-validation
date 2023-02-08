# go-validation

function-chaining-style input validation library for Go and Gin web framework

## Example usages

```go
    package main

    import (
        "net/http"

        "github.com/gin-gonic/gin"

        "example-project/validation"
        "example-project/repository"
    )

    func CreateCustomerUser(c *gin.Context) {
        v := validation.New(c)
        var (
            firstname   = v.Form("firstname").Required().String()
            lastname    = v.Form("lastname").Required().String()
            phoneNumber = v.Form("phone_number").Required().String()
            countryCode = v.Form("country_code").Required().String()
            email       = v.Form("email").Required().String()
            password    = v.Form("password").Required().String()
        )
        if err := v.Done(); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        newUser, err := repository.CreateUser(
            firstname,
            lastname,
            phoneNumber,
            countryCode,
            email,
            password,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to create new user",
            })
            return
        }

        c.JSON(http.StatusCreated, gin.H{
            "result": newUser,
        })
    }

    func main() {
        r := gin.New()

        v1 := r.Group("/api/v1")
        {
            v1.POST("/users", CreateCustomerUser)
        }

        r.Run(":9000")
    }
#
```
