package main

import (
  "encoding/json"
  "fmt"
  "github.com/julienschmidt/httprouter"
  "github.com/streadway/amqp"
  // "github.com/pquerna/ffjson/ffjson"
  "io/ioutil"
  "log"
  "net/http"
)

type session struct {
	*amqp.Connection
	*amqp.Channel
  amqp.Queue
}

var rabbit = GetRabbit()

func GetRabbit() (s session){
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
  defer conn.Close()
  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()
  q, err := ch.QueueDeclare(
    "honeyqa_log_queue", // name
    true, // durable
    false, // delete when unused
    false, // exclusive
    false, // no-wait
    nil, // arguments
  )
  failOnError(err, "Failed to declare a queue")
  return session{conn, ch, q}
}

func InsertLog(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  body, _ := ioutil.ReadAll(r.Body)
  var i map[string]interface{}
  json.Unmarshal(body, &i)
  fmt.Fprintf(w, i["test"].(string))
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}

func main() {
    router := httprouter.New()
    router.POST("/log/insert", InsertLog)
    log.Fatal(http.ListenAndServe(":8080", router))
}
