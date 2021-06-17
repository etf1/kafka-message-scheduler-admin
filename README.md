# kafka-message-scheduler-admin

try it: 

- cd client

- docker run -p 8080:9000 fkarakas/kafka-message-scheduler-admin:mini

- docker build -t ui:beta .    

- docker run -p3000:5000 ui:beta

- open browser on http://localhost:3000


clear backend image : docker rmi -f fkarakas/kafka-message-scheduler-admin:mini



TODO : 
- Afficher icone de tri
- afficher total
- ajouter breadcrumb partout
- refaire complètement écran de détail d’un schedule (avec les versions)
- ajouter des blocs d'info général dans la home
- garder en mémoire la dernière saisie (au moins lors d’un back depuis le breadcrumb)