package main

import (
    "fmt"
    fb "github.com/huandu/facebook"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "time"
)

func main() {

    var access_token = "EAAC9h4ZB7BUABAEqEDmxJZAEVT1LxyUlVdn2Wa1r1WDlMaANtURoOOK0XnsGo9aKuif8HdFUaL2xuoZAZC0VsWGXXsiCA1soWxpeJi8Q0K7jUJ6rzoz2j1SxAgYfFUdaO9rs4ZCOp1w3xvZBbpGUkZAQR5GjTcqhVIZD"
    //var path = "/694472636"
    var items []fb.Result
    var path = "/580099968808112/feed"
    const longForm = "2006-01-02T15:04:05+0000"

    res, _ := fb.Get(path, fb.Params{
        "limit": 25,
        "access_token": access_token,
    })
    //fmt.Printf("res object: %v", res)
    //fmt.Println("Here is my facebook first name:", res["first_name"])

    err := res.DecodeField("data", &items)

    if err != nil {
        fmt.Printf("An error has happened %v", err)
        return
    }

    db, err := sql.Open("mysql", "mysql:mysql@/gorabber?charset=utf8")
    checkErr(err)

    for _, item := range items {
        //fmt.Println(i, item["id"])

        stmt, err := db.Prepare("INSERT posts SET group_id='580099968808112',fb_id=?,updated_at=?,message_text=?")
        checkErr(err)

        var timeStr = fmt.Sprintf("%v",item["updated_time"]);
        updatedAt, _ := time.Parse(longForm, timeStr)

        var messageStr = fmt.Sprintf("%v",item["message"]);
        res, err := stmt.Exec(item["id"], updatedAt, messageStr)
        checkErr(err)

        id, err := res.LastInsertId()
        checkErr(err)
        fmt.Println(item["message"]);
        fmt.Printf("Raw with ID %d inserted!\n", id)

        //updated_time, id
    }





    //fmt.Printf("%v", res)
    //fmt.Print(errno)

    //fmt.Println(access_token)


    /*    res, err := fb.Get("/me/feed", fb.Params{
        "access_token": access_token,
    })

    if err != nil {
        // err can be an facebook API error.
        // if so, the Error struct contains error details.
        if e, ok := err.(*Error); ok {
            fmt.Printf("facebook error. [message:%v] [type:%v] [code:%v] [subcode:%v]",
                e.Message, e.Type, e.Code, e.ErrorSubcode)
            return
        }

        return
    }

    // read my last feed.
    fmt.Println("my latest feed story is:", res.Get("data.0.story"))*/

/*    var first_name string
    res.DecodeField("first_name", &first_name)
    fmt.Println("alternative way to get first_name:", first_name)

    // It's also possible to decode the whole result into a predefined struct.
    type User struct {
        FirstName string
    }

    var user User
    res.Decode(&user)
    fmt.Println("print first_name in struct:", user.FirstName)*/


    }


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}


/*        id
        group_id
        fb_id
        updated_at
        message_text*/