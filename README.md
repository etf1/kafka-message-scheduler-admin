# kafka-message-scheduler-admin

try it: 

- cd client

- docker run -p 8080:9000 fkarakas/kafka-message-scheduler-admin:mini

- docker build -t ui:beta .    

- docker run -p3000:5000 ui:beta

- open browser on http://localhost:3000


clear backend image : docker rmi -f fkarakas/kafka-message-scheduler-admin:mini
