# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicios 1 y 2:
Comando para ejecutar el script:

`./generar-compose.sh <nombre-de-archivo> <cantidad-de-clientes>`

### Ejercicio 3
Comando para iniciar el ambiente de desarrollo:

`make docker-compose-up`

Luego, para validar el echo server:

`./validar-echo-server.sh`

## Ejercicio 4
Comando para iniciar el ambiente de desarrollo:

`make docker-compose-up`

Luego, se puede ejecutar `make docker-compose-down` que a su vez ejecuta `docker compose -f docker-compose-dev.yaml stop -t 1` lo cual envía la señal `SIGTERM` y espera 1 segundo hasta enviar la señal `SIGKILL`


