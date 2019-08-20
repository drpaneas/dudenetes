workflow "Build and deploy on push" {
  on = "push"
  resolves = ["Setup Go for use with actions"]
}

action "Setup Go for use with actions" {
  uses = "actions/setup-go@632d18fc920ce2926be9c976a5465e1854adc7bd"
  runs = "go run dudenetes"
}
