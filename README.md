# kafka-message-scheduler-admin

try it: 

- cd client

- docker run -p 8080:8080 fkarakas/scheduler-admin:beta

- docker build -t ui:beta .    

- docker run -p3000:5000 ui:beta

- open browser on http://localhost:3000


clear backend image : docker rmi -f fkarakas/scheduler-admin:beta
