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


## Ejercicio 6
### Establecimiento de Conexión
Una vez establecida la conexión TCP, el cliente envía un mensaje inicial con el siguiente formato:

```
CANT_BYTES|CLIENT_ID|ACTION
```

donde:
- `CANT_BYTES`: Longitud total del mensaje en caracteres.
- `CLIENT_ID`: Identificación única del cliente.
- `ACTION`: Acción a ejecutar, que en este caso va a ser `BETS` para indicar que enviará apuestas en batches.

### Envío de Apuestas en Batch
Luego, el cliente comienza a enviar apuestas en batch con el siguiente formato:

```
MESSAGE_LENGTH|BATCH
```

donde:
- `MESSAGE_LENGTH`: Longitud total del batch en caracteres.
- `BATCH`: Contiene múltiples apuestas separadas por `BET||BET||BET`.

Cada apuesta dentro del batch sigue el formato:

```
NAME|SURNAME|ID|BIRTHDATE|BET_NUMBER
```

donde:
- `NAME`: Nombre del apostador.
- `SURNAME`: Apellido del apostador.
- `ID`: Identificación del apostador.
- `BIRTHDATE`: Fecha de nacimiento.
- `BET_NUMBER`: Número de apuesta o selección del usuario.

## Respuesta del Servidor
El servidor procesa cada batch recibido y responde con:

- `ACK` → Si el batch es válido y procesado correctamente.
- `ERR` → Si hubo un error en la recepción o procesamiento del batch.

El cliente continúa enviando batches hasta que finaliza la transmisión, o en caso de recibir un ERR deja de enviar batches

## Finalización de la Conexión
Cuando el cliente termina de enviar todas las apuestas, envía el siguiente mensaje para finalizar la conexión:

```
CANT_BYTES|EOF
```

donde:
- `CANT_BYTES`: Longitud total del mensaje en caracteres.
- `EOF`: Indicador de finalización de la transmisión.

El servidor recibe este mensaje y cierra la conexión.

## Ejercicio 7: Consulta de Ganadores
A partir de esta versión del protocolo, el cliente puede enviar otra acción llamada `WINNERS`. Esta request se realiza desde una **nueva conexión TCP** con el siguiente formato:

```
CANT_BYTES|CLIENT_ID|WINNERS
```

donde:
- `CANT_BYTES`: Longitud total del mensaje en caracteres.
- `CLIENT_ID`: Identificación única del cliente.
- `WINNERS`: Acción que solicita los ganadores del sorteo.

### Respuesta del Servidor
El servidor responde dependiendo del estado del sorteo:

- **Si el sorteo se ha realizado y hay ganadores**, la respuesta será:

  ```
  CANT_BYTES|DNI_WINNER1|DNI_WINNER2
  ```
  donde `DNI_WINNER1`, `DNI_WINNER2`, etc., representan los DNI de los ganadores.

- **Si el sorteo se ha realizado pero no hay ganadores en la agencia del cliente**, la respuesta será:

  ```
  CANT_BYTES|NONE
  ```

- **Si el sorteo aún no se ha realizado**, la respuesta será:

  ```
  CANT_BYTES|ERR
  ```

Al recibir la respuesta, la conexión se cierra automáticamente.


## Ejercicio 8: Concurrencia
Se implementó **concurrencia** utilizando la librería `multiprocessing` de Python.

### Implementación
Se creó una nueva clase `ClientHandler` con el objetivo de que cada proceso tenga su propio shutdown. La clase `ClientHandler` incluye un método `shutdown` para cerrar los file descriptors del proceso.


### Herramientas de sincronización del servidor
- **`clients_ready`**: Es un diccionario compartido en el que cada vez que un cliente finaliza el envío de sus batches, se agrega su `client_id`. Esto permite rastrear qué clientes han terminado de enviar sus apuestas.
- **`lottery_completed`**: Es un flag booleano compartido que indica si el sorteo fue realizado.
- **`lottery_lock`**: Se usa para garantizar una sección crítica y evitar que más de un `ClientHandler` anuncie el sorteo.
- **`file_lock`**: Se utiliza para evitar race conditions en el archivo del servidor.

### Flujo del Sorteo
1. Cada cliente envía sus apuestas en batches.
2. Cuando un cliente termina de enviar sus apuestas, se agrega su `client_id` a `clients_ready`.
3. Cuando el último cliente finaliza, el `ClientHandler` correspondiente cambia `lottery_completed` a `True`, permitiendo que futuras solicitudes de `WINNERS` obtengan los resultados.
4. Se usa `lottery_lock` para garantizar que solo un proceso realice este cambio de estado.
