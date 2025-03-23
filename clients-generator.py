import sys, yaml

def generate_compose(num_clients):
    compose = {
        "name": "tp0",
        "services": {
            "server": {
                "container_name": "server",
                "image": "server:latest",
                "entrypoint": "python3 /main.py",
                "environment": [
                    "PYTHONUNBUFFERED=1"
                ],
                "networks": ["testing_net"],
                "volumes": [
                {
                    "type": "bind",
                    "source": "./server/config.ini",
                    "target": "/config/server_config.ini"
                }
                ]
            }
        },
        "networks": {
            "testing_net": {
                "ipam": {
                    "driver": "default",
                    "config": [{"subnet": "172.25.125.0/24"}]
                }
            }
        }
    }

    for i in range(1, num_clients + 1):
        compose["services"][f"client{i}"] = {
            "container_name": f"client{i}",
            "image": "client:latest",
            "entrypoint": "/client",
            "environment": [
                f"CLI_ID={i}",
            ],
            "networks": ["testing_net"],
            "volumes": [
                {
                    "type": "bind",
                    "source": "./client/config.yaml",
                    "target": "/config/client_config.yaml"
                },
                {
                    "type": "bind",
                    "source": "./.data/dataset",
                    "target": "/.data"
                }
                ],
            "depends_on": ["server"]
        }

    return yaml.dump(compose, sort_keys=False, default_flow_style=False)


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
    if num_clients < 0:
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