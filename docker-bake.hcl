variable "GO_VERSION" {
  default = null
}

variable "DESTDIR" {
  default = "./bin"
}

# GITHUB_REF is the actual ref that triggers the workflow and used as version
# when tag is pushed: https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
variable "GITHUB_REF" {
  default = ""
}

target "_common" {
  args = {
    GO_VERSION = GO_VERSION
    GIT_REF = GITHUB_REF
  }
}

group "default" {
  targets = ["binary"]
}

target "binary" {
  inherits = ["_common"]
  target = "binary"
  output = ["${DESTDIR}/build"]
  platforms = ["local"]
}

group "validate" {
  targets = ["lint", "validate-vendor"]
}

target "lint" {
  inherits = ["_common"]
  target = "lint"
  output = ["type=cacheonly"]
}

target "validate-vendor" {
  inherits = ["_common"]
  target = "vendor-validate"
  output = ["type=cacheonly"]
}

target "vendor" {
  inherits = ["_common"]
  target = "vendor-update"
  output = ["."]
}
