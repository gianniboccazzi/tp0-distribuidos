import sys

def generate_compose(num_clients):
    compose_content = """name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net
"""

    for i in range(1, num_clients + 1):
        compose_content += f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={i}
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server
"""

    compose_content += """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

    return compose_content


def parse_arguments():
    try:
        num_clients = int(sys.argv[2])
        output_file = sys.argv[1]
    except ValueError:
        print("Invalid number of clients")
        sys.exit(1)
    return num_clients,output_file


def main():
    num_clients, output_file = parse_arguments()
    if num_clients < 1:
        print("Invalid number of clients")
        sys.exit(1)
    
    compose_content = generate_compose(num_clients)
    with open(output_file, 'w') as f:
        f.write(compose_content)


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 clients-generator.py <output_file> <num_clients>")
        sys.exit(1)
    main()