# Golang Web API Template

Making this so when I need a quick golang server to be setup I can do it in a few commands.

Currently Air is used for the live reloading of the server. This needs to be installed seperately but I need to write a script that can do it automatically.

Same for the migrate I am using https://github.com/golang-migrate/migrate this needs to be installed for the database migrations to happen.

As for the folder structure

├───cmd // This folder contains the main code for the webserver  
│ ├───api // Contains all the routes and the logic for the APIs  
│ └───migrate // Contains the migrations and the migration history for the database  
│ └───migrations  
├───internal // This is where I am storing all the utility functions  
│ ├───db  
│ ├───env  
│ └───store  
└───scripts // I am gonna save all my scripts that needs to be executed for the setup of the server
