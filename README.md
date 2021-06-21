# kafka-message-scheduler-admin

try it: 

- cd client

- docker run -p 8080:9000 fkarakas/kafka-message-scheduler-admin:mini

- docker build -t ui:beta .    

- docker run -p3000:5000 ui:beta

- open browser on http://localhost:3000


clear backend image : docker rmi -f fkarakas/kafka-message-scheduler-admin:mini



TODO : 
OK - Afficher icone de tri
OK - afficher total
OK - ajouter breadcrumb partout
OK - refaire complètement écran de détail d’un schedule (avec les versions)

- ajouter des blocs d'info général dans la home
OK - garder en mémoire la dernière saisie (au moins lors d’un back depuis le breadcrumb)
OK - a jouter un loader au chargement des tables
OK - retailler les colonnes host et al
OK - vider les champs date dans les planif actives
OK - faire un bouton refresh (surtout utile pour les planif actives)
OK - text du nb : 150 plannifications affichées, 316 plannifications au total.
OK - augmenter la taille de la fenetre Valeur en largeur
OK - virer la colonne ordonnanceur
OK - ajouter taille du champ valeur dans la colonne du tableau