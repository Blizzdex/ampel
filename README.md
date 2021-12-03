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


