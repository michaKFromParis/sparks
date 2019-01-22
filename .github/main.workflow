workflow "Build and Publish" {
  on = "push"
  resolves = "Publish"
}

action "Format" {
  uses = "./.github/actions/go"
  args = "format"
}

action "Test" {
  uses = "./.github/actions/go"
  args = "test"
}

action "Build" {
  needs = ["Format", "Test"]
  uses = "./.github/actions/go"
  secrets = ["DOCKER_IMAGE"]
  args = "build"
}

action "Publish Filter" {
  needs = ["Build"]
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "Docker Login" {
  needs = ["Publish Filter"]
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Docker Publish" {
  needs = ["Docker Login"]
  uses = "./.github/actions/docker"
  secrets = ["DOCKER_IMAGE"]
  args = ["publish", "Dockerfile"]
}

action "Publish" {
  needs = ["Publish Filter"]
  secrets = ["GITHUB_TOKEN"]
  uses = "docker://goreleaser/goreleaser:v0.97"
  args = "release"
  }
