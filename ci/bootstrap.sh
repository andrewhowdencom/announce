# Script to bootstrap an environment (e.g. Google Jules) to something that can be used for development.
# Safeties
set -euo pipefail

# Install Taskfile
curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.deb.sh' | sudo -E bash
