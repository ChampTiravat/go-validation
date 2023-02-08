# go-validation

function-chaining-style input validation library for Go and Gin web framework

## Validation options

```go
var (
    firstname   = v.Form("firstname").Required().String()
    Lastname    = v.Param("lastname").Optional().String()
    age         = v.Query("age").Required().Int()
    Height      = v.FormData("height").Required().Float32()
    profilePic  = v.Multipart("profile_picture").Required().File()
)

// String formating (coming soon):
var (
    email    = v.Form("email").Required().String().Format("email")
    userID   = v.Form("user_id").Required().String().Format("numeric")
    postcode = v.Form("postcode").Required().String().Format("numeric").Min(10).Max(20)
)

```

## Example usages

```go
    package main

    import (
        "net/http"

        "github.com/gin-gonic/gin"

        "example-project/validation"
        "example-project/repository"
    )

    func CreateUser(c *gin.Context) {
        v := validation.FromRequest(c)
        var (
            firstname   = v.FormData("firstname").Required().String()
            lastname    = v.FormData("lastname").Required().String()
            phoneNumber = v.FormData("phone_number").Required().String()
            countryCode = v.FormData("country_code").Required().String()
            email       = v.FormData("email").Required().String()
            password    = v.FormData("password").Required().String()
            profilePic  = v.Multipart("profile_picture").Optional().File()
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
            profilePic.,
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
            v1.POST("/users", CreateUser)
        }

        r.Run(":9000")
    }
#
```
