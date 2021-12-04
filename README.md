# Ampel project for the study room at ETH
Contains an updated version of the Ampel, adapted to the current Vis  
Techstack.

# Repo-overview
In this section we shortly go over the content of the repo.
* migrations: Folder containing DB configuration.
* servis: Contains .proto files for GRPC.
* src: Frontend HTML code.
* test: Client to test GRPC requests.
* .gitignore: The .gitginore files specifies the files that are ignored by git.
* Dockerfile: Dockerfile, describes how to create the docker image.
* Makefile: Makefile to start local instance and to compile the .proto files.
* README: this file, contains general information about the repository.
* cinit.yml: Specifies how to launch the ampel inside the docker container.
* docker-compose.yml: Tells docker-compose how to launch a local instance of the ampel.
* go.mod/go.sum: Go dependencies related stuff.
* http.go: Contains the handlers for http requests to the ampel.
* jwt.go: Contains code to extract claims from oidc-jwt-tokens.
* main.go: Contains server setup + grpc handlers of the ampel.
* renovate.json: ???
* sip.yml: Config telling the VIS/SIP infrastructure how to run the ampel.

# How to Run Locally
Make sure you have Docker and Docker-compose installed  
(You should be able to install those via apt, brew or your other favourite package manager)

Open a terminal window in the project repository and execute:
```
sudo docker-compose up --build
```
You "should" now have a local instance of the Ampel running on your computer now.
To check it out just look up 
```
localhost:8080
``` 
on your browser :) You should now see the Ampel homescreen in your browser.

# How does one use it (very difficult)
In order to change the color of the ampel you can go to: 
```
localhost:8080/set
``` 
PS: There really is no no more functionaltiy to it

# A Tour of the project
This section aims at guiding you through the project showing you the most important parts, in order to help you to get an overview.
## Outline
This section aims at leading you through what happens in the ampel when you want to display the ampel color. It starts by looking at the Go code that gets executed and then zooms out to infrastructure running around it.

## Step 1
If you open up main.go you'll find this code in its main function (where execution starts):
```go
//handle http requests
l.Println("Listening")
http.HandleFunc("/set", serv.setColor)
http.HandleFunc("/", serv.getColor)
l.Fatal(http.ListenAndServe(":80", nil))
``` 
What this code does is it tells the application to start listening to HTTP requests on port 80. Each time someone visits ```{URL}/``` or ```{URL}/set```, the two according functions mentioned above are called.
We'll go take a look at serv.getColor (setColor is implemented in a similar way)

## Step 2 (Get & Display the ampel color)
In getColor (located in http.go) we want to request the current color of the ampel and display it.  
This happens in two steps, first we request the ampel color, then we display it.

### 2.1 (Get the ampel color)
We see that getting the color happens in line 
```go
var res, err = s.DbGetColor()
```
s.DbGetColor does a nonstatic call to DbGetColor on our server object. (in main.go) 
```go
sqlStatement := `SELECT color FROM color`
var color int
var err = s.db.QueryRow(sqlStatement).Scan(&color)
if color < 1 || color > 3 {
	log.Warn("Failed to get valid AmpelColor.")
	return 0, err
}
return color, err
```
The only thing it does is sending an SQL query to our postgres DB, that stores the ampels color (yep thats right it stores one color and nothing else...)  
But how did we even set up a connection to any Db anywhere?

### 2.1+ (Set up Db connection)
We've set this connection up when starting the program in main when calling the connectDB function
```go
var dbp, err = sql.Open("postgres", fmt.Sprintf("postgres://%v", *postgresURL))
```
How we got our postgresURL does not really matter for now :)

### 2.2 (Display ampel color)
Now that we have our color, we need to display it on the ampel web interface. The code for this is the second part of the getColor function
```go
//and print the colour to the website.
var p = col4Temp{Col: color}

//create the template if that has not been done yet.
if s.t == nil {
    var e error
    s.t, e = template.ParseFiles("src/colTemplate.html")
    if e != nil {
        l.Fatalf("Failed to parse Template")
    }
}

s.t.Execute(w, p)
return
```
This code creates a template out of our colTemplate.html file (templates behave more nicely than when just serving the html file) and in the template, the .Col variable is not set yet. This variable is set to the ampel color in the second to last line of the above code and displayed online (that's what the Execute command does)  
Now that we saw what is going on inside ampel when displaying its color we can look at what it is running on.


# Step 3 (Do you know Docker)
The Golang code you see in this repository runs encapsulated in a Docker container. This is a somehow independent environment where we can run code in. The Dockerfile contains the blueprints explaining how to build this very container.  
