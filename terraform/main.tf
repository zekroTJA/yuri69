terraform {
  required_providers {
    coder = {
      source  = "coder/coder"
      version = "0.4.1"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 2.16.0"
    }
  }
}

# Admin parameters

# Comment this out if you are specifying a different docker
# host on the "docker" provider below.
variable "step1_docker_host_warning" {
  description = <<-EOF
  This template will use the Docker socket present on
  the Coder host, which is not necessarily your local machine.

  You can specify a different host in the template file and
  surpress this warning.
  EOF
  validation {
    condition     = contains(["Continue using /var/run/docker.sock on the Coder host"], var.step1_docker_host_warning)
    error_message = "Cancelling template create."
  }

  sensitive = true
}

variable "step2_arch" {
  description = <<-EOF
  arch: What architecture is your Docker host on?

  note: codercom/enterprise-* images are only built for amd64
  EOF

  validation {
    condition     = contains(["amd64", "arm64", "armv7"], var.step2_arch)
    error_message = "Value must be amd64, arm64, or armv7."
  }
  sensitive = true
}

variable "step3_OS" {
  description = <<-EOF
  What operating system is your Coder host on?
  EOF

  validation {
    condition     = contains(["MacOS", "Windows", "Linux"], var.step3_OS)
    error_message = "Value must be MacOS, Windows, or Linux."
  }
  sensitive = true
}

variable "workspace_base_image" {
  description = "Which base Docker image would you like to use for your workspace?"
  # The codercom/enterprise-* images are only built for amd64
  default = "codercom/enterprise-base:ubuntu"
  validation {
    condition     = contains(
      ["codercom/enterprise-base:ubuntu", "codercom/enterprise-node:ubuntu", "codercom/enterprise-intellij:ubuntu"], 
      var.workspace_base_image)
    error_message = "Invalid Docker image!"
  }
}

variable "dotfiles_uri" {
  description = <<-EOF
  Dotfiles repo URI (optional)

  see https://dotfiles.github.io
  EOF
    # The codercom/enterprise-* images are only built for amd64
  default = ""
}

# --------------------------------------------------------------------------------------------------

provider "docker" {
  host = var.step3_OS == "Windows" ? "npipe:////.//pipe//docker_engine" : "unix:///var/run/docker.sock"
}

provider "coder" {
}

data "coder_workspace" "me" {
}

resource "coder_agent" "dev" {
  arch = var.step2_arch
  os   = "linux"
  startup_script = var.dotfiles_uri != "" ? "coder dotfiles -y ${var.dotfiles_uri}" : null
}

resource "docker_image" "workspace_image" {
  name = "coder-base-${data.coder_workspace.me.owner}-${lower(data.coder_workspace.me.name)}"
  build {
    path       = "."
    dockerfile = "Dockerfile"
    tag        = ["coder-base-yuri-workspace-image:latest"]
    build_arg = {
      BASE_IMAGE: var.workspace_base_image
    }
  }
}

resource "docker_volume" "home_volume" {
  name = "coder-${data.coder_workspace.me.owner}-${data.coder_workspace.me.name}-root"
}

resource "docker_network" "internal_network" {
  name = "coder-internal-${data.coder_workspace.me.owner}-${data.coder_workspace.me.name}"
  driver = "bridge"
}

resource "docker_container" "workspace" {
  count = data.coder_workspace.me.start_count
  image = docker_image.workspace_image.name
  # Uses lower() to avoid Docker restriction on container names.
  name = "coder-${data.coder_workspace.me.owner}-${lower(data.coder_workspace.me.name)}"
  # Hostname makes the shell more user friendly: coder@my-workspace:~$
  hostname = lower(data.coder_workspace.me.name)
  dns      = ["1.1.1.1"]
  # Use the docker gateway if the access URL is 127.0.0.1
  command = ["bash", "-c", replace(coder_agent.dev.init_script, "127.0.0.1", "host.docker.internal")]
  env     = ["CODER_AGENT_TOKEN=${coder_agent.dev.token}"]
  host {
    host = "host.docker.internal"
    ip   = "host-gateway"
  }
  volumes {
    container_path = "/home/coder/"
    volume_name    = docker_volume.home_volume.name
    read_only      = false
  }
  volumes {
    container_path = "/var/run/docker.sock"
    host_path = "/var/run/docker.sock"
  }
  networks_advanced {
    name = docker_network.internal_network.name
  }
}

resource "docker_volume" "postgres_volume" {
  name = "coder-postgres-${data.coder_workspace.me.owner}-${data.coder_workspace.me.name}"
}

resource "docker_volume" "minio_volume" {
  name = "coder-minio-${data.coder_workspace.me.owner}-${data.coder_workspace.me.name}"
}

resource "docker_container" "lavalink" {
  name = "coder-lavalink-${data.coder_workspace.me.owner}-${lower(data.coder_workspace.me.name)}"
  count = 1
  image = "ghcr.io/zekrotja/lavalink-preconfigured:latest"
  hostname = "lavalink"
  networks_advanced {
    name = docker_network.internal_network.name
  }
}

resource "docker_container" "postgres" {
  name = "coder-postgres-${data.coder_workspace.me.owner}-${lower(data.coder_workspace.me.name)}"
  count = 1
  image = "postgres:alpine"
  hostname = "postgres"
  networks_advanced {
    name = docker_network.internal_network.name
  }
  volumes {
    container_path = "/var/lib/postgresql/data"
    volume_name    = docker_volume.postgres_volume.name
    read_only      = false
  }
  env = [
    "POSTGRES_PASSWORD=yuri69",
    "POSTGRES_USER=yuri69",
    "POSTGRES_DB=yuri69"
  ]
}

resource "docker_container" "minio" {
  name = "coder-minio-${data.coder_workspace.me.owner}-${lower(data.coder_workspace.me.name)}"
  count = 1
  image = "minio/minio:latest"
  hostname = "minio"
  networks_advanced {
    name = docker_network.internal_network.name
  }
  volumes {
    container_path = "/data"
    volume_name    = docker_volume.minio_volume.name
    read_only      = false
  }
  env = [
    "MINIO_ROOT_USER=yuri69",
    "MINIO_ROOT_PASSWORD=yuri69_secret_key"
  ]
  command = ["server", "/data"]
}
