# secret-santa-microservice

The service accepts the POST parameter "persons" with a JSON file.

Example client: https://play.golang.org/p/m84pIyKTDKH

Example JSON input:
---json
[
  {"name":"Name1","email":"example1@gmail.com"},
  {"name":"Name2","email":"example2@gmail.com"}
]
---

Output:
If all ok: ""
If not ok: "{"error":"something wrong"}"
