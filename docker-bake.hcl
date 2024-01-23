variable "GO_VERSION" {
  default = null
}

target "_common" {
  args = {
    GO_VERSION = GO_VERSION
  }
}

group "default" {
  targets = ["vendor"]
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
