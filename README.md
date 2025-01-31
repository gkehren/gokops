# Homelab Infrastructure

This repository contains the Infrastructure as Code (IaC) for managing a homelab server using Ansible and Docker.

## Repository Structure
```
ansible-homelab/
├── ansible.cfg                 # Ansible configuration
├── inventory/                  # Inventory files
│   ├── group_vars/            # Group variables
│   │   └── all/
│   │       ├── vars.yml       # Global variables
│   │       └── vault.yml      # Encrypted sensitive variables
│   ├── host_vars/             # Host-specific variables
│   └── hosts                  # Inventory file
├── playbooks/                 # Task playbooks
│   ├── main.yml              # Main playbook
│   ├── monitoring.yml        # Monitoring stack playbook
│   ├── security.yml          # Security hardening playbook
|   └── applications.yml      # Applications deployment playbook
├── roles/                    # Ansible roles
│   ├── common/              # Common system configurations
│   ├── docker/              # Docker installation and configuration
│   ├── monitoring/          # Prometheus, Grafana, etc.
│   ├── security/           # Security hardening
│   └── applications/       # Applications deployment
└── README.md               # Main documentation
```

## Prerequisites

- Ansible 2.9 or higher
- SSH access to the target server
- Python 3.x installed on the target server

## Quick Start

1. Clone this repository
2. Update the inventory file with your server details
3. Configure variables in group_vars/all/vars.yml
4. Run the main playbook:
```bash
ansible-playbook -i inventory/hosts playbooks/main.yml
```

## Roles Description

### common
Base system configuration including:
- System updates
- Basic packages installation
- User management
- SSH configuration

### docker
Docker installation and configuration:
- Docker CE installation
- Docker Compose installation
- Docker network setup
- Basic Docker configuration

### monitoring
Monitoring stack setup:
- Prometheus
- Grafana
- Node Exporter
- Loki
- Promtail
- AlertManager

### security
Security hardening:
- UFW configuration
- fail2ban setup
- SSH hardening
- Basic security policies

### traefik
Reverse proxy setup:
- Traefik installation
- SSL configuration
- Basic routing setup

### applications
Applications deployment:
- Portainer

## Configuration

### Variables
Main variables are stored in `group_vars/all/vars.yml`. Sensitive data should be stored in `vault.yml` using Ansible Vault.

Example vars.yml:
```yaml
# System Configuration
timezone: "UTC"
system_packages:
  - curl
  - vim
  - htop
  - git
  - python3-pip
  - ufw

# Docker Configuration
docker_compose_version: "2.32.4"
docker_users:
  - "{{ create_user }}"
```

## Security

- All sensitive information should be encrypted using Ansible Vault
- SSH key-based authentication is enforced
- UFW is configured to allow only necessary ports
- Regular system updates are automated

## Backup Strategy

Backup configuration is included in the monitoring role:
- Docker volumes backup
- System configuration backup
- Automated backup rotation
