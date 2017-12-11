package main

import (
    "fmt"
    fb "github.com/huandu/facebook"
    _ "github.com/go-sql-driver/mysql"
    "regexp"
    "database/sql"
    "time"
    "strings"
)

/*CONFIGURATION*/

var db, _ = sql.Open("mysql", "mysql:mysql@/gorabber?charset=utf8")
var depth = 25
var accessToken = "EAAC9h4ZB7BUABAEqEDmxJZAEVT1LxyUlVdn2Wa1r1WDlMaANtURoOOK0XnsGo9aKuif8HdFUaL2xuoZAZC0VsWGXXsiCA1soWxpeJi8Q0K7jUJ6rzoz2j1SxAgYfFUdaO9rs4ZCOp1w3xvZBbpGUkZAQR5GjTcqhVIZD"
var groups = [5]string{"580099968808112"}

func main() {
    var items []fb.Result

    for _, groupId := range groups {
        res, _ := fb.Get("/" + groupId + "/feed", fb.Params{
            "limit":        depth,
            "access_token": accessToken,
        })

        err := res.DecodeField("data", &items)

        if err != nil {
            fmt.Println("No new messages found")
            return
        }


        for _, item := range items {
            var isNew bool
            messageStr := prepareMessage( fmt.Sprintf("%v", item["message"]))
            err := db.QueryRow(
                "SELECT EXISTS(SELECT 1 FROM posts WHERE is_actual = 1 AND group_id = ? AND fb_id = ? AND message_text <> ? )  " +
                "|| NOT EXISTS(SELECT 1 FROM posts WHERE is_actual = 1 AND group_id = ? AND fb_id = ?)",
                groupId,
                item["id"],
                messageStr,
                groupId,
                item["id"]).Scan(&isNew)

            checkErr(err)

            if isNew{
                stmt, err := db.Prepare("UPDATE posts SET is_actual=0 WHERE group_id=? AND fb_id =?")
                checkErr(err)
                res, err := stmt.Exec(groupId, item["id"])
                checkErr(err)
                _ = res
                insertItem(groupId, item)
            }

        }
    }

    }


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func insertItem(groupId string, item fb.Result){
    const longForm = "2006-01-02T15:04:05+0000"
    const isActual = true
    stmt, err := db.Prepare("INSERT posts SET group_id=?,fb_id=?,updated_at=?,message_text=?,is_actual=?")
    checkErr(err)

    var timeStr = fmt.Sprintf("%v", item["updated_time"])
    updatedAt, _ := time.Parse(longForm, timeStr)

    var messageStr = prepareMessage( fmt.Sprintf("%v", item["message"]))
    res, err := stmt.Exec(groupId, item["id"], updatedAt, messageStr, isActual)
    checkErr(err)

    id, err := res.LastInsertId()
    checkErr(err)
    fmt.Printf("Raw with ID %d inserted!\n", id)
}


func prepareMessage(str string) string {
    return strings.ToLower(cutEmojis(str))
}

func cutEmojis(str string) string{
    var emojiRx = regexp.MustCompile(`[\x{1F300}-\x{1F9FF}|\x{2600}-\x{26FF}|[\x{1F680}-\x{1F6FF}]`)
    return emojiRx.ReplaceAllString(str, ``)
}

//"Disponible.  🎄🎄☃️🎄☃️🎄Navidad 🎄☃️🎄☃️HOY DISPONIBLE SE RENTA CASAS EN RENTA PUERTO VALLARTA FECHAS DISPONIBLES DEL 7 al 29 DICIEMBRE NOTA FIN DE AÑO YA ESTA RENTADA ENERO 1 al 30 FEBRERO )==TEMPORADA. BAJA $ 1800 por Noche entra asta 6 personas incluyendo niños menores de ====================por semana completa (pueden entrar asta 6 personas )))===========sercas del🌊🏄🏻🏄🏻🏄🏻🌴🏝🏖 fracc sendero de luna estamos Ubicados dentro de puerto vallarta 4 minutos de la central de autobuses a un costado de la universidad la univa Fracc Sendero de Luna es un coto privado con acceso controlado seguridad las 24 horas y camaras de vigilancia Cuenta con areas de uso comun como Alberca CON CAMASTROS chapoteadero palapa areas verdes ☘️🍀🌴☘️🍀 muy sercas=== Aeropuerto === ✈️✈️zona hotelera == = plaza Marina ⛴🚢🛳🚤⚓️⚓️ centros comercias como plaza Galerias = plaza la Isla🏝🏝costco =walmart=casinos 🎰🎰🎲♣️♠️♥️♦️casa con 3 recamaras muy comodas con camas matrimoniales ventiladores de techo y aire acondicionado cuenta con  🌊🌊🌴🏖🌴baños completos Agunas Casas CUENTA CON asador 🥩🥩🥩🥩 area verde uso exclucivo de la casa =====Internet Aparta tu fecha para esta vacaciones.  ☘️🍀🌴🌴🌴🌊🌊🏄🏻🏄🏻 cel 311-267-61-10)))) whatsApp O Inbox Gracias por compartir esta Publicacion 🌊🏄🏻🌴🌊"