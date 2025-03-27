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


## Ejercicio 5
### Formato de Mensajes
El protocolo de comunicación utiliza mensajes con el siguiente formato:

```
cant_payload|Payload
```

donde:
- `cant_payload`: Indica la longitud del `Payload` en número de caracteres.
- `Payload`: Contiene la información específica del mensaje.

### Mensaje de Apuesta (BET)
Para enviar una apuesta, el mensaje debe seguir la siguiente estructura:

```
MESSAGE_LENGTH|AGENCY|NAME|SURNAME|ID|BIRTHDATE|BET_NUMBER
```

donde:
- `MESSAGE_LENGTH`: Longitud total del payload en caracteres.
- `AGENCY`: Código de la agencia de apuestas.
- `NAME`: Nombre del apostador.
- `SURNAME`: Apellido del apostador.
- `ID`: Identificación del apostador.
- `BIRTHDATE`: Fecha de nacimiento del apostador (formato a definir, ej. `YYYY-MM-DD`).
- `BET_NUMBER`: Número de apuesta o selección del usuario.

### Respuesta del Servidor
El servidor procesa el mensaje y responde con uno de los siguientes valores:

- `ACK` → Si la apuesta se recibió y procesó correctamente.
- `ERR` → Si ocurrió un error en la recepción o procesamiento.

El cliente debe permanecer a la espera de la respuesta después de enviar una apuesta.

