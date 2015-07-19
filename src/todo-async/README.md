Simple pattern for using a json file as a database assuming you don't have too many users or too much data.


Blog Post: [Using a JSON File as a Database Safely in Go](http://txt.fliglio.com/2014/12/safe-json-file-db-in-go/)

## Build

	go get
	go build

## Run

	./jsondb-go

## postgres SQL

DROP SEQUENCE IF EXISTS "public"."tasks_id_seq";
CREATE SEQUENCE "public"."tasks_id_seq" INCREMENT 1 START 25 MAXVALUE 9223372036854775807 MINVALUE 1 CACHE 1;
ALTER TABLE "public"."tasks_id_seq" OWNER TO "postgres";

DROP TABLE IF EXISTS "public"."tasks";
CREATE TABLE "public"."tasks" (
	"id" int4 NOT NULL DEFAULT nextval('tasks_id_seq'::regclass),
	"code" char(36),
	"description" text NOT NULL COLLATE "default"
)
WITH (OIDS=FALSE);
ALTER TABLE "public"."tasks" OWNER TO "postgres";



## Usage

	$ curl -i -H "Content-Type: application/json" -X POST localhost:8080/todo -d '{"description": "Hello World"}'
	HTTP/1.1 201 Created
	Content-Type: application/json
	Date: Sat, 13 Dec 2014 00:11:24 GMT
	Content-Length: 68

	{"id":"5473fc07-314b-4bff-8321-adfdee6decfb","description":"Hello World"}
	
	$ curl -i -H "Content-Type: application/json" -X POST localhost:8080/todo -d '{"description": "Hello World2"}'
	HTTP/1.1 201 Created
	Content-Type: application/json
	Date: Sat, 13 Dec 2014 00:11:35 GMT
	Content-Length: 69

	{"id":"3e39df85-9851-4ce9-af0c-0dd831e3b970","description":"Hello World2"}
	
	$ curl -i -H "Content-Type: application/json" localhost:8080/todo
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Sat, 13 Dec 2014 00:11:48 GMT
	Content-Length: 139

	[{"id":"3e39df85-9851-4ce9-af0c-0dd831e3b970","description”:"Hello World2"},{"id":"5473fc07-314b-4bff-8321-adfdee6decfb","description":"Hello World"}]
	
	$ curl -i -H "Content-Type: application/json" localhost:8080/todo/1
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Sat, 13 Dec 2014 00:12:03 GMT
	Content-Length: 68

	{"id":"5473fc07-314b-4bff-8321-adfdee6decfb","description":"Hello World"}

	$ curl -i -H "Content-Type: application/json" -X DELETE localhost:8080/todo/2
	HTTP/1.1 204 No Content
	Content-Type: application/json
	Date: Sat, 13 Dec 2014 00:12:32 GMT

	$ curl -i -H "Content-Type: application/json" -X PUT localhost:8080/todo/3e39df85-9851-4ce9-af0c-0dd831e3b970 -d '{"description”: "Foo"}'
	HTTP/1.1 200 OK
	Content-Type: application/json
	Date: Sat, 13 Dec 2014 00:13:07 GMT
	Content-Length: 66

	{"id":"3e39df85-9851-4ce9-af0c-0dd831e3b970","description”:"Foo"}

